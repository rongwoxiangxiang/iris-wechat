package main

import (
	"iris/config"
	"iris/routers"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	routers.Routes(app)
	config.InitConfig(app)
	config.Run()
}