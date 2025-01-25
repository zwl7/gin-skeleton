package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	Client *redis.Client
)

// Init 初始化连接
func Init() (err error) {

	fmt.Println(fmt.Sprintf("%d:%s:%s:%d", viper.GetInt("redis.db"), viper.GetString("redis.password"), viper.GetString("redis.host"), viper.GetInt("redis.port")))
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.db"),          // use default DB
		//默认情况下， go-redis 连接池大小为 runtime.GOMAXPROCS * 10，在大多数情况下默认值已经足够使用，且设置太大的连接池几乎没有什么用，可以在 配置项 中调整连接池数量：
		//PoolSize:     viper.GetInt("redis.poolsize"),
		//MinIdleConns: viper.GetInt("redis.minIdleconns"),
	})
	return nil
}

func GetKey(ctx context.Context, key string) (bool, string) {
	val, err := Client.Get(ctx, key).Result()
	if err == nil {
		return true, val
	} else {
		return false, val
	}
}
