package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"net/http"
	"strings"
)

var Register register

type register struct {
}

func (r *register) Register() {
	fmt.Println("开始注册集群。。。")
	client := &http.Client{}

	cluster := new(struct {
		ClusterName string `json:"cluster_name"`
		Ipaddr      string `json:"ipaddr"`
		Port        string `json:"port"`
		K8sVersion  string `json:"k8s_version"`
	})
	//获取k8s集群版本
	nodes, err := Node.GetNodes("", "9999999", "1") //这里写死表示获取所有的node
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, node := range nodes.Items {
		if strings.Contains(node.Name, "master") {
			cluster.K8sVersion = node.Status.NodeInfo.KubeletVersion
		}
	}

	//先从配置文件获取所需参数
	gconf, err := Conf.ReadConfFunc("agent")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cluster.ClusterName = gconf["cluster_name"]
	cluster.Ipaddr = gconf["agent_addr"]
	cluster.Port = gconf["port"]
	// 将结构体编码为 JSON
	jsonData, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("编码结构体为 JSON 时出错：", err)
		return
	}

	// 创建一个包含 JSON 数据的 io.Reader
	jsonReader := bytes.NewReader(jsonData)

	// 创建 HTTP 请求
	//获取server端参数
	sconf, err := Conf.ReadConfFunc("server")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req, err := http.NewRequest(http.MethodPost, "http://"+sconf["server_addr"]+":"+sconf["port"]+"/api/register", jsonReader)
	if err != nil {
		fmt.Println("创建 HTTP 请求报错：", err.Error())
	}

	// 发送 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送 HTTP 请求报错：", err.Error())
	}
	fmt.Println("状态信息：", resp)

	if resp.Status != "200 OK" {
		//开始发送心跳包
		fmt.Println("集群注册失败。。。")
		return
	}
	fmt.Println("开始发送心跳包")
	sendKeepalive(sconf["server_addr"]+":"+sconf["port"], cluster.ClusterName)
	go func() {
		cronSend(sconf["server_addr"]+":"+sconf["port"], cluster.ClusterName)
	}()
	defer resp.Body.Close()
}

var c *cron.Cron

func cronSend(ip string, name string) {

	c = cron.New()
	err := c.AddFunc("*/120 * * * * *", func() {
		sendKeepalive(ip, name)
	})
	if err != nil {
		fmt.Println(err)
	}
	//启动/关闭
	c.Start()
}

func (r *register) CloseCron() {
	c.Stop()
}

// 发送心跳包
func sendKeepalive(ipp string, name string) {
	// 创建 HTTP 客户端
	client := &http.Client{}
	// 创建 HTTP 请求
	req, err := http.NewRequest(http.MethodPost, "http://"+ipp+"/api/keepalive?clusterName="+name, nil)
	if err != nil {
		fmt.Println("创建 HTTP 请求报错：", err.Error())
	}
	// 发送 HTTP 请求
	var resp *http.Response

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("发送 HTTP 请求报错：", err.Error())
	}

	fmt.Println("状态信息：", resp.Status)
	if resp.Status == "200 OK" {
		fmt.Println("保持心跳。。。")
	} else {
		fmt.Println("心跳断开。。。")
	}
	defer resp.Body.Close()
}
