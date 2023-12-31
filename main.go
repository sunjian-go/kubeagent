package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/controller"
	"main/middleware"
	"main/service"
	"net/http"
)

func main() {
	//初始化K8s
	err := service.K8s.K8sInit()
	if err != nil {
		fmt.Println(err.Error())
	}

	//创建路由引擎
	router := gin.Default()
	router.Use(middleware.Cors())
	//初始化路由
	controller.Router.RouterInit(router)

	//启动websocket
	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		http.ListenAndServe(":8082", nil)
		fmt.Println("ws服务已启动。。。")
	}()
	// 运行 Gin 服务
	router.Run(":8081")
}
