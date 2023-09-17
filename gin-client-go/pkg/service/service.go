package service

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"lostar.com/m/pkg/client"
)

func GetServiceList(namespace string) (*v1.ServiceList, error) {
	ctx := context.Background()
	clientSet, err := client.GetK8sClientSet()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	serviceList, err := clientSet.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return serviceList, err
}
