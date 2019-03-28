package routers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"iris/controllers"
)

func Routes(app *iris.Application) {
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Hello id")
	})
	mvc.New(app).Handle(new(controllers.UserController))
}