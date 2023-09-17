package client

import (
	"flag"
	"path/filepath"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

var oneClient sync.Once
var oneConfig sync.Once
var KubeConfig *rest.Config
var KubeClientSet *kubernetes.Clientset

// var kubeConfig *string

// func init() {
// 	// Define kubeconfig
// 	if home := homedir.HomeDir(); home != "" {
// 		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "")
// 	} else {
// 		klog.Fatal("read config error,config is empty")
// 		return
// 	}
// 	flag.Parse()
// }

func GetK8sClientSet() (*kubernetes.Clientset, error) {
	// config, err := GetRestConfig()
	// if err != nil {
	// 	return nil, err
	// }
	// clientSet, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	klog.Fatal(err)
	// 	return nil, err
	// }
	// return clientSet, nil
	oneClient.Do(func() {
		config, err := GetRestConfig()
		if err != nil {
			klog.Fatal(err)
			return
		}
		KubeClientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			klog.Fatal(err)
			return
		}
	})
	return KubeClientSet, nil

}

func GetRestConfig() (config *rest.Config, err error) {
	// config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	// if err != nil {
	// 	klog.Fatal(err)
	// 	return
	// }
	// return config, err
	oneConfig.Do(func() {
		// Define kubeconfig
		var configPath *string
		if home := homedir.HomeDir(); home != "" {
			configPath = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "")
		} else {
			klog.Fatal("read config error,config is empty")
			return
		}
		flag.Parse()
		KubeConfig, err = clientcmd.BuildConfigFromFlags("", *configPath)
		if err != nil {
			klog.Fatal(err)
			return
		}
	})
	return KubeConfig, nil
}
