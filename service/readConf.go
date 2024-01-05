package service

import (
	"github.com/Unknwon/goconfig"
	"github.com/wonderivan/logger"
	"os"
)

var Conf conf

type conf struct {
}

// 读取整个session配置
func (c *conf) ReadConfFunc() (map[string]string, error) {
	currentPath, _ := os.Getwd()
	confPath := currentPath + "/conf/kube_conf.ini"
	_, err := os.Stat(confPath)
	if err != nil {
		logger.Error("file is not found %s")
		return nil, err
	}
	// 加载配置
	config, err := goconfig.LoadConfigFile(confPath)
	if err != nil {
		logger.Error("读取配置文件出错:", err)
		return nil, err
	}
	// 获取 section
	kubeconf, _ := config.GetSection("kube")
	//fmt.Println("配置文件内容：", kubeconf)
	//fmt.Println("websocker地址：", kubeconf["wshost"])
	return kubeconf, nil
}
