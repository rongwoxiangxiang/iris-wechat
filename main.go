package main

import (
	"github.com/kataras/iris"
	"iris/config"
	"iris/routers"
)

func main() {
	app := iris.New()
	routers.Routes(app)
	config.InitConfig(app)
	config.Run()
}