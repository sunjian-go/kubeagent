package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/controller"
	"main/middle"
	"main/service"
)

func main() {
	//初始化K8s
	err := service.K8s.K8sInit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//创建路由引擎
	router := gin.Default()
	router.Use(middle.Cors())
	//router.Use(middle.JWTAuth()) //加载jwt中间件，用于token验证
	//初始化路由
	controller.Router.RouterInit(router)

	//开一个携程去注册集群
	go func() {
		service.Register.Register()
	}()

	// 运行 Gin 服务
	router.Run(":8081")
	//关闭cron
	defer service.Register.CloseCron()
}
