package service

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"main/conf"
)

var K8s k8s

type k8s struct {
	Clientset *kubernetes.Clientset `json:"clientset"`
	Conf      *rest.Config          `json:"conf"`
}

func (k *k8s) K8sInit() error {
	// 配置 Kubernetes 客户端
	config, err := clientcmd.BuildConfigFromFlags("", conf.Kubeconfig)
	if err != nil {
		return fmt.Errorf("获取 Kubernetes 配置失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("创建 Kubernetes 客户端失败: %v", err)
	}
	k.Clientset = clientset
	k.Conf = config
	return nil
}
