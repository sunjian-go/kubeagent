package service

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"mime/multipart"
)

var CopyTpod copyTpod

type copyTpod struct {
}

type PodInfo struct {
	PodName       string `form:"podName"`
	Namespace     string `form:"namespace"`
	ContainerName string `form:"containerName"`
	Path          string `form:"path"`
}

// CopyToPod 将文件保存到 Kubernetes Pod 容器中
func (c *copyTpod) CopyToPod(files []*multipart.FileHeader, podinfo *PodInfo) error {
	// 确保至少有一个文件
	if len(files) == 0 {
		fmt.Println("没有上传文件")
		return errors.New("没有上传文件")
	}

	// 创建一个 tar 归档的 buffer
	var tarBuffer bytes.Buffer      //创建缓冲区：bytes.Buffer是一个实现了io.ReadWriter接口的类型，可以作为数据缓冲区
	tw := tar.NewWriter(&tarBuffer) //tar.NewWriter函数用于创建一个新的Tar文件写入器，会将Tar文件的内容写入到提供的io.Writer接口实现中（这里将tarBuffer作为io.Writer接口传递给了tar.NewWriter，因此会向tarBuffer中写内容）
	defer tw.Close()                //确保函数执行完毕关闭该文件写入器

	// 将所有文件添加到 tar 归档中,有几个文件就循环几次
	for _, file := range files {
		src, err := file.Open() //返回一个 io.Reader 接口，用于读取每个上传文件的内容
		if err != nil {
			return fmt.Errorf("打开文件失败: %v", err)
		}
		defer src.Close()

		// 创建 tar 头部信息: 目的是使用Tar文件写入器将文件内容写入到Tar文件时，每个文件条目都会包含这些元数据信息，使得在解压或读取Tar文件时能够正确地还原文件的属性。
		hdr := &tar.Header{
			Name: file.Filename, //文件名
			Mode: 0600,          //只读
			Size: file.Size,     //文件大小
		}

		//将hdr所代表的文件头信息写入到Tar文件中
		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("写入 tar 头部失败: %v", err)
		}

		//拷贝文件到tar文件写入器,最终写入到 tarBuffer 中
		if _, err := io.Copy(tw, src); err != nil {
			return fmt.Errorf("写入 tar 归档失败: %v", err)
		}
	}

	// 确保所有内容都写入了 tar 归档 (关闭写入器以确保所有数据都被正确地刷新到文件系统。)
	if err := tw.Close(); err != nil {
		return fmt.Errorf("关闭 tar 写入器失败: %v", err)
	}

	//构建一个可以在指定 Pod 内部执行 tar 命令的 POST 请求对象
	command := []string{"tar", "xvf", "-", "-C", podinfo.Path} //创建字符串切片存储命令，’-‘代表从标准输入读取
	req := K8s.Clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podinfo.PodName).
		Namespace(podinfo.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
			Container: podinfo.ContainerName,
		}, scheme.ParameterCodec)

	//创建SPDY 执行器（SPDY Executor）来创建一个可以在pod内远程执行命令的执行器对象（可以将上面定义的请求对象发送到 Kubernetes API 服务器，并在指定的 Pod 内部执行相应的命令。）
	exec, err := remotecommand.NewSPDYExecutor(K8s.Conf, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("创建执行器失败: %v", err)
	}

	// 执行远程命令,将 tar 归档传输到 Pod内
	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  bytes.NewReader(tarBuffer.Bytes()), //使用 bytes.NewReader(tarBuffer.Bytes()) 创建一个 bytes.Reader 对象，将之前创建的 tarBuffer 的内容作为标准输入流传递给命令。
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return fmt.Errorf("执行失败: %v", err)
	}

	return nil
}
