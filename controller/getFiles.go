package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
	"net/http"
)

var File file

type file struct {
}

func (f *file) GetFilesForWeb(c *gin.Context) {
	podinfo := new(service.PodInfo)
	podinfo.PodName = c.Query("podName")
	podinfo.Namespace = c.Query("namespace")
	podinfo.ContainerName = c.Query("containerName")
	podinfo.Path = c.Query("path")
	//fmt.Println("aaaaaaaaaaaaaaaaa", podinfo)
	if podinfo.PodName == "" || podinfo.Namespace == "" || podinfo.Path == "" {
		fmt.Println("pod信息不完善，请设置完再上传")
		c.JSON(400, gin.H{
			"err": "pod信息不完善，请设置完再上传",
		})
		return
	}
	//podinfo.PodName = "nginx"
	//podinfo.Namespace = "sjtest"
	// 从请求中获取文件
	formFiles, err := c.MultipartForm() //将请求体中的数据解析为一个 multipart.Form 对象，然后，通过访问 formFiles.File 字段，你可以获取表单中指定字段名的文件列表。
	if err != nil {
		c.String(http.StatusBadRequest, "获取上传文件失败: %v", err)
		return
	}
	files := formFiles.File["file"] // 这里的 "file" 是前端表单里的字段名

	//进行拷贝文件
	err = service.CopyTpod.CopyToPod(files, podinfo)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "上传文件到pod成功！",
	})
}

func (f *file) DownLoadFile(c *gin.Context) {
	podinfo := new(service.PodInfo)
	if err := c.ShouldBindJSON(podinfo); err != nil {
		c.JSON(400, gin.H{
			"err": "数据绑定失败：" + err.Error(),
		})
		return
	}
	fmt.Println("获取数据为：", podinfo)
	err := service.CopyTpod.CopyFromPod(podinfo, c)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.Status(200)
}
