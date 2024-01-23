package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
)

var Node node

type node struct {
}

func (n *node) GetNodes(c *gin.Context) {
	//GET请求
	nodes, err := service.Node.GetNodes()
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
