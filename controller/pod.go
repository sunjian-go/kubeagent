package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"main/service"
)

var Pod pod

type pod struct {
}

// 此结构体用于内部，用来绑定客户端传过来的pod信息
type podInfo struct {
	FilterName string `form:"filter_name"`
	NameSpace  string `form:"namespace"`
	Limit      int    `form:"limit"`
	Page       int    `form:"page"`
}

type logPod struct {
	Container string `form:"container"`
	Podname   string `form:"podname"`
	Namespace string `form:"namespace"`
}

// 获取pod列表，支持分页、过滤、排序
func (p *pod) GetPods(c *gin.Context) {
	pod := new(podInfo)
	if err := c.Bind(pod); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定pod参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	fmt.Println("客户端传过来的为：", *pod)
	podlist, err := service.Pod.GetPods(pod.FilterName, pod.NameSpace, pod.Limit, pod.Page)
	if err != nil {
		logger.Info("获取pod列表失败, " + err.Error())
		c.JSON(400, gin.H{
			"err":  "获取pod列表失败, " + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取列表成功",
		"data": *podlist,
	})
}

// 获取容器信息
func (p *pod) GetContainer(c *gin.Context) {
	pod := new(service.PodDetail)
	if err := c.Bind(pod); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败" + err.Error(),
			"data": nil,
		})
		return
	}
	containers, err := service.Pod.GetContainer(pod.Name, pod.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  "获取容器信息失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取容器信息成功",
		"data": containers,
	})
}

// 获取容器日志
func (p *pod) GetContainerLog(c *gin.Context) {
	pod := new(logPod)
	if err := c.Bind(pod); err != nil {
		c.JSON(400, gin.H{
			"err": "数据绑定失败" + err.Error(),
		})
		return
	}
	err := service.Pod.GetPodLog(pod.Container, pod.Podname, pod.Namespace, c)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  "获取日志失败：" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取日志成功",
		"data": nil,
	})
}
