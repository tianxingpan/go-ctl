// Package cmd provides parsing configuration files
package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// 设置go-ctl配置或环境变量

// 设置环境变量
func SetEnv(env string) error {
	arr := strings.Split(env, "=")
	if len(arr) < 2 {
		return fmt.Errorf("")
	}
	_ = os.Setenv(arr[0], arr[1])
	return nil
}

// 设置配置
func SetConfig(key string, value interface{}) {
	viper.Set(key, value)
}

