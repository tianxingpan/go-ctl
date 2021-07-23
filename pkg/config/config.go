// Package cmd provides parsing configuration files
package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
)

// 默认配置定义
var (
	defaultConfigPath = ".go-ctl"
	defaultConfigName = "config"
	defaultConfigType = "yaml"
)

func Init(conf string) error {
	return initConfig(conf)
}

// 生成默认配置文件
func generateDefaultConf() {
	_, err := os.Stat(defaultConfigPath)
	isCreate := false
	// 判断目录是否存在，不存在则创建
	if os.IsNotExist(err) {
		err = os.Mkdir(defaultConfigPath, 0777)
		if err != nil {
			fmt.Printf("\033[0;32;31m")
			fmt.Printf("mkdir failed: %s\n", err)
			fmt.Printf("\033[m")
		}
		isCreate = true
	}
	// 判断文件是否存在
	fn := fmt.Sprintf("%s.%s", defaultConfigName, defaultConfigType)
	fPath := path.Join(defaultConfigPath, fn)
	_, err = os.Stat(fPath)
	if os.IsNotExist(err) {
		isCreate = true
	}
	if isCreate {
		// 创建配置文件
		f, err := os.Create(fPath)
		if err != nil {
			fmt.Printf("\033[0;32;31m")
			fmt.Printf("%s\n", err)
			fmt.Printf("\033[m")
			return
		}
		f.Close()
	}
}

// 初始化viper
func initConfig(conf string) error {
	// 获取home目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("\033[0;32;31m")
		fmt.Printf("Unable to get home directory: %s\n", err)
		fmt.Printf("\033[m")
		return err
	}
	defaultConfigPath = path.Join(homeDir, defaultConfigPath)
	// 如果指定了配置文件，则解析指定的配置文件
	// 如果没有指定配置文件，则解析默认的配置文件
	if conf != "" {
		viper.SetConfigFile(conf)
	} else {
		generateDefaultConf()
		viper.AddConfigPath(defaultConfigPath)
		viper.SetConfigName(defaultConfigName)
	}

	// 设置配置文件格式为YAML
	viper.SetConfigType(defaultConfigType)

	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("\033[0;32;31m")
		fmt.Printf("viper read config failed: %s\n", err)
		fmt.Printf("\033[m")
		return err
	}
	return nil
}

// 刷新默认配置
func RefreshConfig(s bool) {
	if s {
		// 刷新配置
		fn := fmt.Sprintf("%s.%s", defaultConfigName, defaultConfigType)
		filepath := path.Join(defaultConfigPath, fn)
		err := viper.WriteConfigAs(filepath) // 直接写入，有内容就覆盖，没有文件就新建
		if err != nil {
			fmt.Printf("\033[0;32;31m")
			fmt.Printf("refresh config failed: %s\n", err)
			fmt.Printf("\033[m")
		}
	}
}