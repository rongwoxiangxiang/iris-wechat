package config

import (
	"fmt"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"time"
)
type RedisConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Pass string `yaml:"pass"`
	Database string `yaml:"database"`
	MaxIdle int `yaml:"maxIdle"`
	MaxActive int `yaml:"maxActive"`
	IdleTimeout time.Duration `yaml:"timeout"`
	Prefix string `yaml:"prefix"`
}


var (
	redisClient *redis.Database
)

func CacheSet(key string, value interface{}, expire time.Duration)  {
	GetRedisClient().Set("wechat",sessions.LifeTime{Time: time.Now().Add(expire * time.Second)} , key, value, false)
}

func CacheGet(key string) (value interface{}) {
	return GetRedisClient().Get("wechat", key)
}

func GetRedisClient() *redis.Database {
	if redisClient == nil {
		configs := GetConfigs()
		if configs != nil {
			Addr := fmt.Sprintf("%s:%s",configs.RedisConfig.Host , configs.RedisConfig.Port)
			redisClient = redis.New(service.Config{
				Network:     "tcp",
				Addr:        Addr,
				Password:    configs.RedisConfig.Pass,
				Database:    configs.RedisConfig.Database,
				MaxIdle:     configs.RedisConfig.MaxIdle,
				MaxActive:   configs.RedisConfig.MaxActive,
				IdleTimeout: configs.RedisConfig.IdleTimeout,
				Prefix:      configs.RedisConfig.Prefix})
		}
	}
	return redisClient
}
