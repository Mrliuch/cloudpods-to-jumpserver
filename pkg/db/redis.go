package db

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Rdb *redis.Client

func init() {
	redisHost := viper.GetString("redis.host")
	redisPort := viper.GetString("redis.port")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	logrus.Infof("Redis地址：%s", redisAddr)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // 密码
		DB:       0,  // 数据库
		PoolSize: 20, // 连接池大小
	})
	Rdb = rdb
}
