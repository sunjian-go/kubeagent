package service

import (
	"github.com/Unknwon/goconfig"
	"main/utils"
	"os"
)

var Conf conf

type conf struct {
}

// 读取整个session配置
func (c *conf) ReadConfFunc(opt string) (map[string]string, error) {
	currentPath, _ := os.Getwd()
	confPath := currentPath + "/conf/conf.ini"
	_, err := os.Stat(confPath)
	if err != nil {
		utils.Logg.Error("配置文件未找到")
		return nil, err
	}
	// 加载配置
	config, err := goconfig.LoadConfigFile(confPath)
	if err != nil {
		utils.Logg.Error("读取配置文件出错:" + err.Error())
		return nil, err
	}
	// 获取 section
	var gconf map[string]string
	switch opt {
	case "server":
		gconf, _ = config.GetSection("server")
		break
	case "agent":
		gconf, _ = config.GetSection("agent")
		break

	}
	return gconf, nil
}
