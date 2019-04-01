package config

import (
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"time"
)

var (
	redisClient *redis.Database
)

func CacheSet(key string, value interface{}, expire time.Duration)  {
	redisClient.Set("wechat",sessions.LifeTime{Time: time.Now().Add(expire * time.Second)} , key, value, false)
}

func CacheGet(key string) (value interface{}) {
	return redisClient.Get("wechat", key)
}

func GetRedisClient() *redis.Database {
	return redisClient
}

func init()  {
	if redisClient == nil {
		redisClient = redis.New(service.Config{
			Network:     "tcp",
			Addr:        "127.0.0.1:6379",
			Password:    "",
			Database:    "",
			MaxIdle:     0,
			MaxActive:   10,
			IdleTimeout: service.DefaultRedisIdleTimeout,
			Prefix:      ""})
	}
}
