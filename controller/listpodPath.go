package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
)

var Listpath listpath

type listpath struct {
}

func (l *listpath) ListContainerPath(c *gin.Context) {
	podinfo := new(service.PodInfo)
	if err := c.ShouldBindJSON(podinfo); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败" + err.Error(),
		})
		return
	}
	fmt.Println("pod信息：", podinfo)
	out, err := service.Listpath.ListContainerPath(podinfo)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "列出容器路径成功",
		"data": out,
	})
}
