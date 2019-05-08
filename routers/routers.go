package routers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"iris/controllers"
	"time"
)

func Routes(app *iris.Application) {
	app.Get("/check", func(ctx iris.Context) {
		ctx.WriteString( time.Now().Format(time.RFC3339))
	})
	app.PartyFunc("/service/{flag:string}", func(r iris.Party) {
		mvc.New(r).Handle(new(controllers.ServiceController))
	})

}