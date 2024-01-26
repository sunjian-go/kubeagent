package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
)

var Pack packet

type packet struct {
}

// 启动抓包进程
func (p *packet) StartPacket(c *gin.Context) {
	url := c.Query("url")
	packinfo := new(service.PackInfo)
	if err := c.ShouldBindJSON(packinfo); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	fmt.Println("需要抓包的数据为：", packinfo)
	err := service.Pack.StartPacket(packinfo, url)
	if err != nil {
		if err.Error() == "当前已有抓包程序运行，请稍后重试" {
			c.JSON(442, gin.H{
				"err": nil,
			})
		} else {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
		}
		return
	}
	c.JSON(200, gin.H{
		"msg": "启动抓包程序成功",
	})
}

func (p *packet) StopPacket(c *gin.Context) {
	url := c.Query("url")
	err := service.Pack.StopPacket(c, url)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "关闭抓包程序成功",
	})
}
