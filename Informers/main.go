package main

import (
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientSet, 30*time.Second)
	deployInformer := informerFactory.Apps().V1().Deployments()
	informer := deployInformer.Informer()
	deployLister := deployInformer.Lister()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	stopper := make(chan struct{})
	defer close(stopper)
	// 启动Informer List & Watch
	informerFactory.Start(stopper)
	// 等待所有的Informer 缓存同步
	informerFactory.WaitForCacheSync(stopper)
	deployments, err := deployLister.Deployments("default").List(labels.Everything())
	// 遍历deploy 列表
	for index, deploy := range deployments {
		fmt.Printf("%d -> %s\n", index, deploy.Name)
	}

	<-stopper
}

func onAdd(obj interface{}) {
	deploy := obj.(*v1.Deployment)
	klog.Infoln("add a deploy: ", deploy.Name)
}

func onUpdate(old, new interface{}) {
	oldDeploy := old.(*v1.Deployment)
	newDeploy := new.(*v1.Deployment)
	klog.Infoln("update deploy: ", oldDeploy.Name, newDeploy.Name)
}

func onDelete(obj interface{}) {
	deploy := obj.(*v1.Deployment)
	klog.Infoln("delete deploy: ", deploy.Name)
}
