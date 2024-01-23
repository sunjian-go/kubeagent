package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Node node

type node struct {
}
type NodeResp struct {
	Items []corev1.Node `json:"items"`
	Total int           `json:"total"`
}

// 获取node列表
func (n *node) GetNodes() (noderesp *NodeResp, err error) {
	nodeList, err := K8s.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取node列表失败: " + err.Error())
		return nil, errors.New("获取node列表失败: " + err.Error())
	}

	return &NodeResp{
		Items: nodeList.Items,
		Total: len(nodeList.Items),
	}, nil
}
