package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
)

var (
	application *iris.Application
)

func InitConfig(app *iris.Application)  {
	application = app
	initDb(application)
}

func Run()  {
	application.Run(runner(), configuration())
}

func runner() iris.Runner {
	return iris.Addr("localhost:8080")
}
func configuration() iris.Configurator {
	configurate := iris.YAML("config/iris.yml")
	return iris.WithConfiguration(configurate)
}
