package service

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var K8s k8s

type k8s struct {
	Clientset *kubernetes.Clientset `json:"clientset"`
	Conf      *rest.Config          `json:"config"`
}

func (k *k8s) K8sInit() error {
	// 配置 Kubernetes 客户端
	//config, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	config, err := rest.InClusterConfig() //根据集群中的环境变量和配置文件创建一个 *rest.Config 对象，该对象包含与集群通信所需的认证信息和服务器地址等配置。

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
