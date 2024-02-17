package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"time"
)

var Terminal terminal

type terminal struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源的WebSocket连接
		return true
	},
}

// WsHandler是一个WebSocket请求处理函数，用于处理从前端发送来的WebSocket请求
func (t *terminal) WsHandler(namespace, podName, containerName, bashType string, c *gin.Context) error {
	fmt.Println("有客户端连接")

	//创建一个TerminalSession类型的pty实例,用于向websocket读写信息
	pty, err := NewTerminalSession(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("获取pty实例失败: %v\n", err)
		return err
	}
	//处理关闭
	defer func() {
		logger.Info("关闭ws连接")
		pty.Close()
	}()

	// 设置连接关闭的回调函数
	pty.wsConn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("WebSocket 连接关闭，状态码：%d，原因：%s\n", code, text)
		return nil
	})

	fmt.Println("ws已连接。。。")

	//组装post请求，请求内容为执行在容器中的命令
	// 初始化pod所在的corev1资源组
	// PodExecOptions struct 包括Container stdout stdout Command 等结构
	// scheme.ParameterCodec 应该是pod 的GVK （GroupVersion & Kind）之类的
	// URL长相:
	// https://192.168.1.11:6443/api/v1/namespaces/default/pods/nginx-wf2-778d88d7c7rmsk/exec?
	//command=%2Fbin%2Fbash&container=nginxwf2&stderr=true&stdin=true&stdout=true&tty=true
	req := K8s.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{bashType},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	logger.Info("exec post request url: ", req)

	//升级SPDY协议
	executor, err := remotecommand.NewSPDYExecutor(K8s.Conf, "POST", req.URL())
	if err != nil {
		logger.Error("建立SPDY连接失败，" + err.Error())
		return err
	}

	//与kubelet建立stream连接
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
	if err != nil {
		logger.Error("执行pod命令失败，" + err.Error())
		//将报错返回出去
		pty.Write([]byte("执行pod命令失败，" + err.Error()))
		//标记退出stream流
		pty.Done()
	}
	return nil
}

// 消息内容
type terminalMessage struct {
	Operation string `json:"operation"` // 操作类型，比如"stdin" "stdout" "stderr"
	Data      string `json:"data"`      // 数据内容，比如命令内容，命令执行结果等
	Rows      uint16 `json:"rows"`      // 终端的行数，用于resize操作
	Cols      uint16 `json:"cols"`      // 终端的列数，用于resize操作
}

// 交互的结构体，接管输入和输出
type TerminalSession struct {
	wsConn   *websocket.Conn                 // WebSocket连接
	sizeChan chan remotecommand.TerminalSize // 用于传输终端大小的channel
	doneChan chan struct{}                   // 用于标记WebSocket连接关闭的channel
}

// 创建TerminalSession类型的对象并返回
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, errors.New("升级websocket失败：" + err.Error())
	}
	// 创建TerminalSession实例
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return session, nil
}

// 重写Read方法，给内部调用
// Read用于从WebSocket连接中读取消息，接收web端输入的指令内容,返回值int是读成功了多少数据
func (t *TerminalSession) Read(p []byte) (int, error) {
	//设置ws超时时间
	t.wsConn.SetReadDeadline(time.Now().Add(10 * time.Minute))

	//从ws中读取消息,也就是读取stdin的消息
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		log.Printf("读取stdin的消息失败: %v", err)
		return 0, err
	}
	fmt.Println("读取到来自前端的消息：", message)
	//从ws中读取出来的stdin的消息进行反序列化
	var msg terminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("反序列化消息失败: %v", err)
		return 0, err
	}
	//fmt.Println("反序列化之后为：", msg)
	//根据消息内容的选项做不同动作
	switch msg.Operation {
	case "stdin":
		//fmt.Println("msg.data= ", msg.Data)
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		log.Printf("unknown message type '%s'", msg.Operation)
		return 0, fmt.Errorf("unknown message type '%s'", msg.Operation)
	}
}

// 重写write方法
// Write用于拿到apiserver的返回内容，向WebSocket连接中写入消息
func (t *TerminalSession) Write(p []byte) (int, error) {
	//将apiserver的返回内容组装进结构体并进行编码
	msg, err := json.Marshal(terminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("组装消息结构体失败：%v", err)
		return 0, err
	}
	//开始写数据
	//fmt.Println("写入数据：", msg)
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("写入消息失败：%v", err)
		return 0, err
	}
	//返回写入数据的长度
	return len(p), nil
}

// Done用于标记WebSocket连接已经关闭
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

// Close用于关闭WebSocket连接
func (t *TerminalSession) Close() {
	t.wsConn.Close()
}

// Next用于获取下一个终端大小，或者在WebSocket连接关闭时返回nil
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan: //读取到size数据的话就返回该数据
		return &size
	case <-t.doneChan: //读取到数据的话就代表关闭ws
		return nil
	}
}
