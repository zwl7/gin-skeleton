package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	//设置加载配置文件的路径
	viper.SetConfigFile("./conf/config.yaml")

	//读取配置
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("viper.ReadInConfig fail")
		return
	}

	//配置文件修改后，自动加载最新的配置
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改了，已重新加载配置")
	})

	//fmt.Println(viper.GetInt("port"))

	fmt.Println("domain-------")
	fmt.Println(viper.GetString("domain"))

	return
}
