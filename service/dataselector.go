package service

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	Networkingv1 "k8s.io/api/networking/v1"
	"sort"
	"strings"
	"time"
)

// 排序对象
type FilterQuery struct {
	Name string
}

// 分页对象
type PaginateQuery struct {
	Limit int //每页的数量
	Page  int //页数
}

// DataSelectQuery 定义过滤和分页的属性，过滤：Name， 分页：Limit和Page
// Limit是单页的数据条数
// Page是第几页
type DataSelectQuery struct {
	FilterQuery   *FilterQuery   //过滤用
	PaginateQuery *PaginateQuery //分页用
}

// dataSelect 用于封装排序、过滤、分页的数据类型
type DataSelector struct {
	GenericDataList []DataCell       //元素切片
	DataSelectQuery *DataSelectQuery //用于定义过滤和分页的属性的结构体
}

// DataCell接口，用于各种资源list的类型转换，转换后可以使用dataSelector的自定义排序方法
type DataCell interface {
	GetCreation() time.Time //该方法用于获取资源的时间戳
	GetName() string        //该方法用户获取资源的Name
}

// 排序
// 实现自定义结构的排序，需要重写Len、Swap、Less方法
// Len方法用于获取数组长度
func (d *DataSelector) Len() int {
	return len(d.GenericDataList) //计算资源切片的长度返回
}

// Swap方法用于数组中的元素在比较大小后的位置交换，可定义升序或降序
func (d *DataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less方法用于定义数组中元素排序的“大小”的比较方式
func (d *DataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation() //获得a的时间
	b := d.GenericDataList[j].GetCreation() //获得b的时间
	return b.Before(a)                      //如果b在a时间之前则为真
}

// 重写以上3个方法用使用sort.Sort进行排序
func (d *DataSelector) Sort() *DataSelector {
	//排完序直接返回
	sort.Sort(d)
	return d
}

// 过滤
// Filter方法用于过滤元素，比较元素的Name属性，若包含，再返回
func (d *DataSelector) Filter() *DataSelector {
	//fmt.Println("传进来的Name为：", d.dataSelectQuery.FilterQuery.Name)
	//若Name的传参为空，则返回所有元素
	if d.DataSelectQuery.FilterQuery.Name == "" {
		return d
		fmt.Println("Name为空！！！")
	}
	//若Name的传参不为空，则返回元素名中包含Name的所有元素
	filterdList := []DataCell{}
	for _, value := range d.GenericDataList {
		if !strings.Contains(value.GetName(), d.DataSelectQuery.FilterQuery.Name) {
			continue
		} else {
			//fmt.Println("找到了：", value.GetName())
			filterdList = append(filterdList, value)
		}
	}
	//fmt.Println("找到了pods: ", filterdList)
	d.GenericDataList = filterdList
	//fmt.Println("找到了pods: ", d.GenericDataList)
	return d
}

// 分页
// Paginate方法用于数组分页，根据Limit和Page的传参，返回一定范围的数据
func (d *DataSelector) Paginate() *DataSelector {
	//fmt.Println("开始分页，limit为：", d.DataSelectQuery.PaginateQuery.Limit)
	limit := d.DataSelectQuery.PaginateQuery.Limit //获取当前元素切片每页的数量
	page := d.DataSelectQuery.PaginateQuery.Page   //获取当前元素切片的页数
	fmt.Println("分页信息 ", limit, page)
	//验证参数合法，若参数不合法，则返回所有数据
	if limit <= 0 || page <= 0 {
		//fmt.Println("分页失败后：", d.GenericDataList)
		return d
	}
	//定义取数范围的起始索引和结束索引：
	//如果传入的limit数量大于等于deployment列表的长度，就把列表实际长度给limit,这时候page就肯定是1了，因为只有一页
	//fmt.Println("列表长度为：", len(d.GenericDataList))
	if limit >= len(d.GenericDataList) {
		limit = len(d.GenericDataList)
		page = 1
	}
	//起始索引值
	startindex := limit * (page - 1)
	fmt.Println("起始为：", startindex)
	endIndex := 0
	//endIndex不能减1，因为切片默认会减1
	endIndex = limit * page
	//处理最后一页，这时候就把endIndex由30改为25了
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	fmt.Println("起始：", startindex, " 结尾：", endIndex)
	d.GenericDataList = d.GenericDataList[startindex:endIndex] //从下标为0的到下标为最后一个的元素赋值给元素切片
	fmt.Println("成功分页后个数为：", len(d.GenericDataList))
	return d
}

// pod
// 定义podCell类型，实现DataCell接口，用于类型转换
// corev1.Pod->podCell->DataCell，podCell相当于corev1.Pod转到DataCell的一个桥梁
type podCell corev1.Pod

// podCell重写DataCell的GetCreation方法
func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time //返回时间戳
}

// podCell重写DataCell的GetName方法
func (p podCell) GetName() string {
	return p.Name //返回Name
}

// deployment
// 定义deploymentCell类型，实现DataCell接口，用于类型转换
// appsv1.Deployment->deploymentCell->DataCell，deploymentCell相当于appsv1.Deployment转到DataCell的一个桥梁
type deploymentCell appsv1.Deployment

// deploymentCell重写DataCell的GetName方法
func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

// deploymentCell重写DataCell的GetCreation方法
func (d deploymentCell) GetName() string {
	return d.Name
}

// daemonSet
// 定义daemonSetCell类型，实现DataCell接口，用于类型转换
type daemonSetCell appsv1.DaemonSet

// daemonSetCell重写DataCell的GetName方法
func (d daemonSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

// daemonSetCell重写DataCell的GetCreation方法
func (d daemonSetCell) GetName() string {
	return d.Name
}

// statefulSet
// 定义statefulSetCell类型，实现DataCell接口，用于类型转换
type statefulSetCell appsv1.StatefulSet

// statefulSetCell重写DataCell的GetName方法
func (d statefulSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

// statefulSetCell重写DataCell的GetCreation方法
func (d statefulSetCell) GetName() string {
	return d.Name
}

// service
type serviceCell corev1.Service

func (s serviceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}
func (s serviceCell) GetName() string {
	return s.Name
}

// ingress
type ingressCell Networkingv1.Ingress

func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}
func (i ingressCell) GetName() string {
	return i.Name
}

// configMap
type cmCell corev1.ConfigMap

func (c cmCell) GetCreation() time.Time {
	return c.CreationTimestamp.Time
}
func (c cmCell) GetName() string {
	return c.Name
}

// Secret
type secretCell corev1.Secret

func (s secretCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}
func (s secretCell) GetName() string {
	return s.Name
}

// pv
type pvCell corev1.PersistentVolume

func (p pvCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func (p pvCell) GetName() string {
	return p.Name
}

// pvc
type pvcCell corev1.PersistentVolumeClaim

func (p pvcCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func (p pvcCell) GetName() string {
	return p.Name
}

// namespace
type nsCell corev1.Namespace

func (p nsCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func (p nsCell) GetName() string {
	return p.Name
}
