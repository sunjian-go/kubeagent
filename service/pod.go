package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"main/config"
	"main/utils"
	"time"
)

// 定义pod类型和Pod对象，用于包外的调用(包是指service目录)，例如Controller
var Pod pod

type pod struct {
}

type PodDetail struct {
	Name      string `form:"name"`
	Namespace string `form:"namespace"`
}

// 定义列表的返回内容，Items是pod元素列表，Total为pod元素数量
type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}
type PodsNp struct {
	Namespace string `json:"namespace"`
	PodNum    int    `json:"podNum"`
}

// 定义PodsNp类型，用于返回namespace中pod的数量
type PodNp struct {
	Namespace string `json:"namespace"`
	PodNum    int    `json:"podNum"`
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (p *pod) toCell(pods []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

// 获取pod列表
// 获取pod列表，支持过滤、排序、分页
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsresp *PodsResp, err error) {

	//获取podList类型的pod列表
	//context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label，field等
	//kubectl get services --all-namespaces --field-seletor metadata.namespace != default
	podList, err := K8s.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	//fmt.Println("podlist: ", podList.Items)
	if err != nil {
		//logger用于打印日志
		//return用于返回response内容
		utils.Logg.Error("获取pod列表失败" + err.Error())
		return nil, errors.New("获取pod列表失败" + err.Error())
	}
	//实例化dataSelector对象
	selecttableData := &DataSelector{
		GenericDataList: p.toCell(podList.Items), //将pods列表转换为DataCell类型赋值
		DataSelectQuery: &DataSelectQuery{ //
			FilterQuery: &FilterQuery{
				Name: filterName, //将传进来的需要查找的Name赋值给该结构体
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit, //将传进来的页数和每页的数量赋值
				Page:  page,
			},
		},
	}
	//先过滤
	filterd := selecttableData.Filter()
	//total := len(selecttableData.GenericDataList)
	total := len(filterd.GenericDataList) //计算过滤好的目标pod列表的长度
	//fmt.Println("过滤后：", filterd.GenericDataList)
	//再排序和分页
	//for _, data := range filterd.GenericDataList {
	//	fmt.Println("排序分页前：", data.GetName(), data.GetCreation())
	//}
	pods := filterd.Sort().Paginate() //连续调用排序和分页方法
	//for _, data := range pods.GenericDataList {
	//	fmt.Println("排序分页后：", data.GetName(), data.GetCreation())
	//}
	//将[]DataCell类型的pod列表转为v1.pod列表
	podv1s := p.fromCells(pods.GenericDataList)
	//fmt.Println("整理好的pod信息：", podv1s)
	return &PodsResp{
		Items: podv1s,
		Total: total,
	}, nil
}

// 获取pod详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		utils.Logg.Error("获取pod详情失败: " + err.Error())
		return nil, errors.New("获取pod详情失败: " + err.Error())
	}
	return pod, nil
}

// 获取容器信息
func (p *pod) GetContainer(podName, namespace string) (containers []string, err error) {
	//获取pod详情
	pod, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		utils.Logg.Error("获取pod详情失败: " + err.Error())
		return nil, errors.New("获取pod详情失败: " + err.Error())
	}
	//从pod中拿到容器名
	for _, cont := range pod.Spec.Containers {
		containers = append(containers, cont.Name)
	}
	return containers, nil
}

// 获取pod内容器日志
func (p *pod) GetPodLog(containerName, podName, namespace string, c *gin.Context) (err error) {
	utils.Logg.Info("开始获取日志")
	//new一个TerminalSession类型的pty实例,用来向前端发送数据
	pty, err := NewTerminalSession(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logg.Error("get pty failed: " + err.Error())
		return errors.New("get pty failed: " + err.Error())
	}
	//设置超时时间
	pty.wsConn.SetWriteDeadline(time.Now().Add(10 * time.Minute))

	//1.设置日志的配置，容器名，获取的内容的配置
	lineLimit := int64(config.PodLogTailLine) //先将定义的行数转为int64位
	option := &corev1.PodLogOptions{          //定义一个corev1.PodLogOptions指针并赋值
		Container: containerName,
		TailLines: &lineLimit,
		Follow:    true,
	}

	//2.获取一个request实例
	req := K8s.Clientset.CoreV1().Pods(namespace).GetLogs(podName, option)
	//3.发起stream连接，得到Response.body
	podLog, err := req.Stream(context.TODO())
	if err != nil {
		utils.Logg.Error("获取podLog失败" + err.Error())
		return errors.New("获取podLog失败" + err.Error())
	}

	defer func() {
		podLog.Close() //关闭stream连接
		pty.Close()    //关闭pty连接
		utils.Logg.Info("pty连接已关闭")
	}()
	//4.将response.body写入到缓冲区，目的是为了转换成string类型
	buf := make([]byte, 4096)
	//循环读取日志，并通过socket发送给前端
	for {
		size, err := podLog.Read(buf)
		if size > 0 {
			_, err := pty.Write(buf[:size]) //写入前端
			//当报错的时候就是前端关闭了socket连接
			if err != nil {
				utils.Logg.Error("pty写入报错" + err.Error())
				break
			}
			//fmt.Println("获取到日志，开始写入：", string(buf))
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read pod logs error: %v", err)
		}
	}
	//5.转换数据返回
	return nil
}
