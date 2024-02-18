package service

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"main/config"
	"main/utils"
)

var K8s k8s

type k8s struct {
	Clientset *kubernetes.Clientset `json:"clientset"`
	Conf      *rest.Config          `json:"config"`
}

func (k *k8s) K8sInit() error {
	// 配置 Kubernetes 客户端
	gconf, err := Conf.ReadConfFunc("server")
	if err != nil {
		utils.Logg.Error(err.Error())
		return err
	}
	//根据模式选择k8s的初始化方式
	switch gconf["model"] {
	case "docker":
		k.Conf, err = clientcmd.BuildConfigFromFlags("", config.Kubeconfig) //读取config的方式
		break
	case "k8s":
		//rest.InClusterConfig() 函数的作用是获取 Kubernetes 集群的配置信息，以便与集群进行交互。
		//它通过检查环境变量来确定当前程序是否在 Kubernetes 集群内部运行。如果在集群内部运行，它将返回集群的配置信息，
		//包括 API 服务器的地址和证书等。如果不在集群内部运行，它将返回错误。
		k.Conf, err = rest.InClusterConfig() //部署的时候只需给工作负载绑定一个具有集群最高权限的sa即可
		break
	}

	if err != nil {
		return fmt.Errorf("获取 Kubernetes 配置失败: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(k.Conf)
	if err != nil {
		return fmt.Errorf("创建 Kubernetes 客户端失败: %v", err)
	}
	k.Clientset = clientset
	return nil
}
