package client

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

var kubeConfig *string

func init() {
	// Define kubeconfig
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "")
	} else {
		klog.Fatal("read config error,config is empty")
		return
	}
	flag.Parse()
}

func GetK8sClientSet() (*kubernetes.Clientset, error) {
	config, err := GetRestConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
		return nil, err
	}
	return clientSet, nil
}

func GetRestConfig() (config *rest.Config, err error) {
	config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		klog.Fatal(err)
		return
	}
	return config, err
}
