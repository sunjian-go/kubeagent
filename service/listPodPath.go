package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

var Listpath listpath

type listpath struct {
}

func (l *listpath) ListContainerPath(podinfo *PodInfo) (string, error) {

	fmt.Println("pod信息：", podinfo.PodName, podinfo.ContainerName, podinfo.Namespace)
	req := K8s.Clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podinfo.PodName).
		Namespace(podinfo.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   []string{"ls", "/"}, // List the root directory
			Stdout:    true,
			Stderr:    true,
			Container: podinfo.ContainerName,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(K8s.Conf, "POST", req.URL())
	if err != nil {
		fmt.Println("创建执行器失败：" + err.Error())
		return "", errors.New("创建执行器失败：" + err.Error())
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	if err != nil {
		fmt.Println("执行命令失败：" + err.Error())
		return "", errors.New("执行命令失败：" + err.Error())
	} else {
		fmt.Println("列出路径为: \n", stdout.String())
		fmt.Println("报错为:", stderr.String())

		return stdout.String(), nil
	}

}
