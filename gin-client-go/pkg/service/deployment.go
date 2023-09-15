package service

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"lostar.com/m/pkg/client"
)

func GetDeployment(namespace string, name string) (*v1.Deployment, error) {
	ctx := context.Background()
	clientSet, err := client.GetK8sClientSet()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	deployment, err := clientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	return deployment, err
}

func ListDeployment(namespace string) ([]v1.Deployment, error) {
	ctx := context.Background()
	clientSet, err := client.GetK8sClientSet()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	deploymentList, err := clientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	return deploymentList.Items, err
}
