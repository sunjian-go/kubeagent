package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
)

var Icmp icmp

type icmp struct {
}

// ping方法
func (i *icmp) PingFunc(c *gin.Context) {
	url := c.Query("url")
	icmp := new(service.Icmpdata)
	if err := c.ShouldBindJSON(icmp); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}

	icmpresp, err := service.Icmp.PingFunc(icmp, url)
	if err != nil {
		if err.Error() == "err" {
			c.JSON(400, gin.H{
				"err": icmpresp,
			})
		} else {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
		}
		return
	}
	c.JSON(200, icmpresp)

}
