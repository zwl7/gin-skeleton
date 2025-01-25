package mysql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

// 👉🏻 https://gorm.io/zh_CN/docs/

var (
	DB *gorm.DB
)

func InitMySQL() error {

	//获得一个*grom.DB对象
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"))

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// 不能直接DB，err := gorm.Open() 因为:=会重新定义新的变量DB，DB就变成了局部变量，外面访问时就变成了nil
	//DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	err := errors.New("")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: viper.GetString("mysql.table_prefix"),
		},
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("Gorm init 异常：", err)
	}

	//根据*grom.DB对象获得*sql.DB的通用数据库接口
	sqlDB, err := DB.DB()

	//max_open_conns
	sqlDB.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns")) //设置最大连接数
	sqlDB.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns")) //设置最大的空闲连接数

	data, err := json.Marshal(sqlDB.Stats()) //获得当前的SQL配置情况
	fmt.Println(string(data))

	return err
}
