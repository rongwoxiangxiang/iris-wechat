package config

import (
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions/sessiondb/redis"
)

var (
	orm *xorm.Engine
)

func initDb(application *iris.Application)  {
	var err error
	orm, err = xorm.NewEngine("mysql", "root:Cs@229229@(cdb-n4dhxbt3.gz.tencentcdb.com:10085)/jwechat?charset=utf8")
	orm.ShowSQL(true)
	if err != nil {
		application.Logger().Fatalf("orm failed to initialized: %v", err)
	}
	redisClient = redis.New() //可选择在redis服务器之间配置网桥
}

func GetDb() *xorm.Engine {
	return orm
}