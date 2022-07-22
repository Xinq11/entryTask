package database

import (
	"EntryTask/config"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var RedisDB *redis.Client

func RedisInit() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		PoolSize: config.RedisPoolNum,
	})
	_, err := RedisDB.Ping().Result()
	if err != nil {
		logrus.Panic("Redis connect error: ", err.Error())
	}
	logrus.Infoln("Redis connect success")
}
