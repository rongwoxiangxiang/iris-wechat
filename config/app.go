package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"time"
)

var (
	Orm *xorm.Engine
	redisClient *redis.Database
	application *iris.Application
)

func InitConfig(app *iris.Application)  {
	application = app
	initDb()
	iris.RegisterOnInterrupt(func() {
		Orm.Close()
		redisClient.Close()
	})
	defer func() {
		Orm.Close()
		redisClient.Close()
	}()
}

func Run()  {
	application.Run(runner(), configuration())
}

func initDb()  {
	var err error
	Orm, err = xorm.NewEngine("mysql", "root:Cs@229229@(cdb-n4dhxbt3.gz.tencentcdb.com:10085)/jwechat?charset=utf8")
	Orm.ShowSQL(true)
	if err != nil {
		application.Logger().Fatalf("orm failed to initialized: %v", err)
	}
	redisClient = redis.New() //可选择在redis服务器之间配置网桥
}

func runner() iris.Runner {
	return iris.Addr("localhost:8080")
}
func configuration() iris.Configurator {
	return iris.WithConfiguration(iris.YAML("config/iris.yml"))
}

func CacheSet(key string, value interface{}, expire time.Duration)  {
	redisClient.Set("wechat",sessions.LifeTime{Time: time.Now().Add(expire * time.Second)} , key, value, false)
}

func CacheGet(key string) (value interface{}) {
	return redisClient.Get("wechat", key)
}


