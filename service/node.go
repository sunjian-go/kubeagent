package service

import (
	"context"
	"errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"main/utils"
	"strconv"
)

var Node node

type node struct {
}
type NodeResp struct {
	Items []corev1.Node `json:"items"`
	Total int           `json:"total"`
}

type NodeInfo struct {
	FilterName string `json:"filter_name"`
	Limit      string `json:"limit"`
	Page       string `json:"page"`
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (n *node) toCell(node []corev1.Node) []DataCell {
	cells := make([]DataCell, len(node))
	for i := range node {
		cells[i] = nodeCell(node[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (n *node) fromCells(cells []DataCell) []corev1.Node {
	node := make([]corev1.Node, len(cells))
	for i := range cells {
		node[i] = corev1.Node(cells[i].(nodeCell))
	}
	return node
}

// 获取node列表
func (n *node) GetNodes(filterName, limit, page string) (noderesp *NodeResp, err error) {
	nodeList, err := K8s.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logg.Error("获取node列表失败: " + err.Error())
		return nil, errors.New("获取node列表失败: " + err.Error())
	}

	//实例化dataSelector对象
	limitnum, _ := strconv.Atoi(limit)
	pagenum, _ := strconv.Atoi(page)
	selecttableData := &DataSelector{
		GenericDataList: n.toCell(nodeList.Items), //将pods列表转换为DataCell类型赋值
		DataSelectQuery: &DataSelectQuery{ //
			FilterQuery: &FilterQuery{
				Name: filterName, //将传进来的需要查找的Name赋值给该结构体
			},
			PaginateQuery: &PaginateQuery{
				Limit: limitnum, //将传进来的页数和每页的数量赋值
				Page:  pagenum,
			},
		},
	}

	filterd := selecttableData.Filter()   //先过滤
	total := len(filterd.GenericDataList) //计算过滤好的目标pod列表的长度
	nodes := filterd.Sort().Paginate()    //连续调用排序和分页方法

	//将[]DataCell类型的pod列表转为v1.pod列表
	nodev1s := n.fromCells(nodes.GenericDataList)
	return &NodeResp{
		Items: nodev1s,
		Total: total,
	}, nil

}

// 获取node详情
func (n *node) GetNodeDetail(nodeName string) (node *corev1.Node, err error) {
	node, err = K8s.Clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		utils.Logg.Error("获取node详情失败: " + err.Error())
		return nil, errors.New("获取node详情失败: " + err.Error())
	}
	return node, nil
}
