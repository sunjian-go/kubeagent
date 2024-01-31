package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
)

var Node node

type node struct {
}

// 获取node列表
func (n *node) GetNodes(c *gin.Context) {
	//GET请求
	node := new(service.NodeInfo)
	node.FilterName = c.Query("filter_name")
	node.Limit = c.Query("limit")
	node.Page = c.Query("page")
	fmt.Println("node信息：", node)
	nodes, err := service.Node.GetNodes(node.FilterName, node.Limit, node.Page)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取node列表成功",
		"data": nodes,
	})
}

// 获取node详情
func (n *node) GetNodeDetail(c *gin.Context) {
	nodeName := c.Query("nodeName")
	fmt.Println("nodeName为： ", nodeName)
	nodedetail, err := service.Node.GetNodeDetail(nodeName)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取node详情成功",
		"data": nodedetail,
	})
}
