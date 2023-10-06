package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
	"lostar.com/m/pkg/client"
)

func GetPodList(namespace string) (*v1.PodList, error) {
	ctx := context.Background()

	clientSet, err := client.GetK8sClientSet()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	podList, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return podList, err
}

func GetPod(namespace, name string) (*v1.Pod, error) {
	ctx := context.Background()
	clientSet, err := client.GetK8sClientSet()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	pod, err := clientSet.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return pod, err
}

type WsMessage struct {
	MessageType int
	data        []byte
}

type WsConnection struct {
	wsSocket  *websocket.Conn
	inChan    chan *WsMessage
	outChan   chan *WsMessage
	mutex     sync.Mutex
	isClosed  bool
	classChan chan byte
}

func (wsConn *WsConnection) WsClose() {
	err := wsConn.wsSocket.Close()
	if err != nil {
		klog.Error(err)
		return
	}
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.classChan)
	}
}

func (wsConn *WsConnection) wsReadLoop() {
	var (
		msgType int
		data    []byte
		msg     *WsMessage
		err     error
	)
	for {
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			goto ERROR
		}
		msg = &WsMessage{
			MessageType: msgType,
			data:        data,
		}
		select {
		case wsConn.inChan <- msg:
		case <-wsConn.classChan:
			goto CLOSED
		}
	}
ERROR:
	wsConn.WsClose()
CLOSED:
}
func (wsConn *WsConnection) wsWriteLoop() {
	var (
		msg *WsMessage
		err error
	)
	for {
		select {
		case msg = <-wsConn.outChan:
			if err = wsConn.wsSocket.WriteMessage(msg.MessageType, msg.data); err != nil {
				goto ERROR
			}
		case <-wsConn.classChan:
			goto CLOSED
		}
	}
ERROR:
	wsConn.WsClose()
CLOSED:
}

func (wsConn *WsConnection) WsWrite(MessageType int, data []byte) (err error) {
	select {
	case wsConn.outChan <- &WsMessage{MessageType: MessageType, data: data}:
		return
	case <-wsConn.classChan:
		err = errors.New("websocket closed")
	}
	return
}

func (wsConn *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	case msg = <-wsConn.inChan:
		return
	case <-wsConn.classChan:
		err = errors.New("websocket closed")
	}
	return
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWebsocket(resp http.ResponseWriter, req *http.Request) (wsConn *WsConnection, err error) {
	var (
		wsSocket *websocket.Conn
	)
	if wsSocket, err = wsUpgrader.Upgrade(resp, req, nil); err != nil {
		klog.Error(err)
		return
	}
	wsConn = &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *WsMessage, 1000),
		outChan:   make(chan *WsMessage, 1000),
		classChan: make(chan byte),
		isClosed:  false,
	}
	go wsConn.wsReadLoop()
	go wsConn.wsWriteLoop()
	return
}

type streamHeander struct {
	wsConn      *WsConnection
	resizeEvent chan remotecommand.TerminalSize
}

func (headler *streamHeander) Write(p []byte) (size int, err error) {
	copyData := make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	err = headler.wsConn.WsWrite(websocket.TextMessage, copyData)
	return
}

type xtermMessage struct {
	MsgType string `json:"type"`
	Input   string `json:"input"`
	Rows    uint16 `json:"rows"`
	Cols    uint16 `json:"cols"`
}

func (headler *streamHeander) Read(p []byte) (size int, err error) {
	var xtermMsg xtermMessage
	msg, err := headler.wsConn.WsRead()
	if err != nil {
		klog.Errorln(err)
		return
	}
	if err = json.Unmarshal(msg.data, &xtermMsg); err != nil {
		return
	}
	if xtermMsg.MsgType == "resize" {
		headler.resizeEvent <- remotecommand.TerminalSize{
			Width:  uint16(xtermMsg.Cols),
			Height: uint16(xtermMsg.Rows),
		}
	} else if xtermMsg.MsgType == "input" {
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return
}

func (handler *streamHeander) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.resizeEvent
	size = &ret
	return
}

// func (headler *streamHeander) Close() error {
// 	return nil
// }

// func (headler *streamHeander) Flush() error {
// 	return nil
// }

func WebSSH(namespaceName, podName, containerName, method string, resp http.ResponseWriter, req *http.Request) error {
	var (
		err      error
		executor remotecommand.Executor
		wsConn   *WsConnection
	)
	ctx := context.Background()
	config, err := client.GetRestConfig()
	if err != nil {
		return err
	}
	clientSet, err := client.GetK8sClientSet()
	if err != nil {
		return err
	}
	reqSSH := clientSet.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespaceName).SubResource("exec").
		VersionedParams(
			&v1.PodExecOptions{
				Container: containerName,
				Command:   []string{method},
				Stderr:    true,
				Stdout:    true,
				Stdin:     true,
				TTY:       true,
			},
			scheme.ParameterCodec)
	if executor, err = remotecommand.NewSPDYExecutor(config, "POST", reqSSH.URL()); err != nil {
		klog.Errorln(err)
		return err
	}
	if wsConn, err = InitWebsocket(resp, req); err != nil {
		return err
	}
	handler := &streamHeander{wsConn: wsConn, resizeEvent: make(chan remotecommand.TerminalSize)}
	if err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		goto END
	}
END:
	klog.Errorln(err)
	wsConn.WsClose()
	return err
}
