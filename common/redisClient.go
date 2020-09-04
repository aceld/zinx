package common

import (
	"fmt"
	"wsserver/configs"
	"github.com/go-redis/redis"
)

var (
	client *redis.Client
)

func InitRedis() {
	redisAddr := fmt.Sprintf("%s:%d", configs.GConf.RedisAddr, configs.GConf.RedisPort)
	client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: configs.GConf.RedisPw, // no password set
		DB:       0,                     // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println("初始化redis:", pong, err)
}

func GetRedisClient() (c *redis.Client) {
	return client
}