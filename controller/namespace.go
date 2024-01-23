package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
)

var Namespace namespace

type namespace struct {
}

// 获取namespace列表
func (n *namespace) GetNamespaces(c *gin.Context) {
	namespaces, err := service.Namespace.GetNamespaces()
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	//fmt.Println("获取到namespacelist: ", namespaces)
	c.JSON(200, gin.H{
		"msg":  "获取namespace列表成功",
		"data": namespaces,
	})
}
