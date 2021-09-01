package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	VERSION     string
	PORT        int
	DB_TYPE     string
	DB_USER     string
	DB_HOST     string
	DB_PORT     int
	DB_PASS     string
	DB_NAME     string
	SECRET      string
	ISSUER      string
	STREAM_PING int
}

func Get() *Configuration {
	config := Configuration{}
	viper := viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// 自定义部署设置环境变量
	// LEEKBOX_ENV=custom
	// LEEKBOX_DIR=指定配置根目录.默认$HOME/.leekbox/
	if os.Getenv("LEEKBOX_ENV") == "custom" {
		if os.Getenv("LEEKBOX_HOME") != "" {
			viper.AddConfigPath(os.Getenv("LEEKBOX_HOME"))
		} else {
			viper.AddConfigPath(fmt.Sprintf("%s/.leekbox/", os.Getenv("HOME")))
		}
	}
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("致命错误 未找到配置文件: %s\n", err))
	}
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("致命错误 配置文件格式不对: %s\n", err))
	}
	// viper.WatchConfig()
	// viper.OnConfigChange(func(in fsnotify.Event) {
	// 	fmt.Println("配置文件不能随便修改！")
	// 	if err := viper.Unmarshal(config); err != nil {
	// 		panic(fmt.Errorf("致命错误 配置文件格式不对 %s\n", err))
	// 	}
	// 	fmt.Printf("config: %+v\n", config)
	// })
	return &config
}
