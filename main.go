package main

import (
	"iris/config"
	"iris/routers"
	_"iris/routers"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	routers.Routes(app)
	app.Run(config.Server())
}