package config

import (
	"fmt"
	"github.com/chnyangzhen/kago-fly/pkg/constant"
	"github.com/chnyangzhen/kago-fly/pkg/helper"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	// 显示调用Set设置值 > 命令行参数（flag）> 环境变量-配置文件 > key/value存储 > 默认值
	root *viper.Viper
)

// ViperLifecycle Viper组件生命周期
type ViperLifecycle int

// Configuration Viper配置包装器
type Configuration struct {
	namespace string
	root      *viper.Viper
	alias     map[string]string // 本地Key别名
}

// GetWrapper 获取Viper配置包装器
func GetWrapper(namespace string) *Configuration {
	return &Configuration{
		namespace: namespace,
		root:      root,
		alias:     make(map[string]string),
	}
}

// GlobalConfig 获取全局配置
func GlobalConfig() *viper.Viper {
	IsInitialized()
	return root
}

// IsInitialized 是否已经初始化
func IsInitialized() {
	if root == nil {
		panic("Viper component is not initialized, Please load the viper component ")
	}
}

func InitConfig(configNames []string) error {
	// 获取所有的配置文件
	err := filepath.Walk(constant.DefaultConfigPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dir := filepath.Dir(path)
			viper.AddConfigPath(dir)
			// 读取配置文件
			filename := filepath.Base(path)
			filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename)) // 去除扩展名
			if helper.ContainsString(configNames, filenameWithoutExt) || helper.ContainsString(configNames, filename) {
				viper.SetConfigName(filenameWithoutExt)
				err := viper.MergeInConfig()
				if err != nil {
					return fmt.Errorf("failed to read config file %s: %v", path, err)
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error loading config files: %v：%s\n", configNames, err)
	}
	root = viper.GetViper()
	root.AutomaticEnv()
	return nil
}
