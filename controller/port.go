package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
)

var Port portt

type portt struct {
}

func (p *portt) PortTel(c *gin.Context) {
	url := c.Query("url")
	portdata := new(service.PortData)
	if err := c.ShouldBindJSON(portdata); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	info, err := service.Port.TCPTelnet(portdata, url)
	if err != nil {
		if err.Error() == "err" {
			c.JSON(400, gin.H{
				"err": info,
			})
		} else {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
		}
		return
	}
	c.JSON(200, gin.H{
		"msg": info,
	})
}
