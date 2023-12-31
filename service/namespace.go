package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Namespace namespace

type namespace struct {
}
type namespaceResp struct {
	Namespaces []corev1.Namespace `json:"namespaces"`
	Total      int                `json:"total"`
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (n *namespace) toCell(namespace []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(namespace))
	for i := range namespace {
		cells[i] = nsCell(namespace[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (n *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	ns := make([]corev1.Namespace, len(cells))
	for i := range cells {
		ns[i] = corev1.Namespace(cells[i].(nsCell))
	}
	return ns
}

// 获取namespace列表
func (n *namespace) GetNamespaces() (namespaces *namespaceResp, err error) {
	namespaceList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace列表失败：" + err.Error())
		return nil, errors.New("获取namespace列表失败：" + err.Error())
	}
	return &namespaceResp{
		Namespaces: namespaceList.Items,
		Total:      len(namespaceList.Items),
	}, nil
}
