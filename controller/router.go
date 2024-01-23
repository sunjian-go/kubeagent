package controller

import (
	"github.com/gin-gonic/gin"
)

var Router router

type router struct {
}

func (r *router) RouterInit(router *gin.Engine) {
	router.
		//POST("/api/login", Login.Login).
		POST("/api/upload", File.GetFilesForWeb).
		GET("/api/corev1/getpods", Pod.GetPods).
		GET("/api/corev1/getcontainers", Pod.GetContainer).
		GET("/api/corev1/getnamespaces", Namespace.GetNamespaces).
		GET("/api/corev1/getnodes", Node.GetNodes).
		GET("/api/corev1/getlog", Pod.GetContainerLog).
		GET("/api/terminal", Pod.TerminalFunc).
		GET("/api/listPath", Listpath.ListContainerPath).
		POST("/api/download", File.DownLoadFile)

}
