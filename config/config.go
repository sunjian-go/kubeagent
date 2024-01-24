package config

const (
	//Kubeconfig = "E:\\Code\\project\\kubernetes\\kubecp\\k8sconf\\conf.txt" //本地
	Kubeconfig     = "/root/.kube/config" //正式环境
	PodLogTailLine = 2000                 //tail的日志行数 tail -n 2000
	AAA            = 123
)
