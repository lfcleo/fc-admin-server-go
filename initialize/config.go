package initialize

import (
	"fc-admin-server-go/pkg/config"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func ConfigSetUp() {
	// 也可以使用SetConfigFile直接指定
	viper.SetConfigFile("./config.yaml")

	// 查找并读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		//如果配置文件读取错误，可以直接注销程序了。
		panic(fmt.Errorf("【Panic】启动配置文件错误:: %w", err))
		return
	}

	//将配置文件映射到结构体
	if err := viper.Unmarshal(config.Data); err != nil {
		//如果配置文件读取错误，可以直接注销程序了。
		panic(fmt.Errorf("【Panic】配置文件映射错误:: %w", err))
	}

	//监听修改
	viper.WatchConfig()
	//为监听配置修改增加一个回调函数,当值被修改时，重新赋值
	viper.OnConfigChange(func(in fsnotify.Event) {
		//将配置文件映射到结构体
		if err := viper.Unmarshal(config.Data); err != nil {
			//如果配置文件读取错误，可以直接注销程序了。
			panic(fmt.Errorf("【Panic】配置文件映射错误:: %w", err))
		}
	})
}
