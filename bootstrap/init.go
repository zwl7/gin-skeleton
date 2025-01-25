package bootstrap

import (
	"errors"
	"fmt"
	"gin-skeleton/database/mysql"
	"gin-skeleton/database/redis"
	"gin-skeleton/global/variable"
	"gin-skeleton/logger"
	"gin-skeleton/settings"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {

	err := settings.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	//初始化日志
	if err := logger.Init(); err != nil {
		fmt.Println("logger init error:", err)
		return
	}
	//将缓冲区的日志，追加到文件中
	defer zap.L().Sync()

	//InitMySQL
	if err := mysql.InitMySQL(); err != nil {
		fmt.Println("mysql init error:", err)
		return
	}

	if err := redis.Init(); err != nil {
		fmt.Println("redis init error:", err)
		return
	}

	//初始化全局变量
	variable.ServerDomain = viper.GetString("domain")
	variable.ErrorUserNotLogin = errors.New("用户未登录")

}
