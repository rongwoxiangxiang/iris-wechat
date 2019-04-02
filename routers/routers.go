package routers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"iris/controllers"
)

func Routes(app *iris.Application) {
	app.Get("/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		ctx.Writef("Hello id: %d", id)
	})
	app.PartyFunc("/service/{flag:string}", func(r iris.Party) {
		mvc.New(r).Handle(new(controllers.ServiceController))
	})

}