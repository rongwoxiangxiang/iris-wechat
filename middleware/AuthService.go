package middleware

import (
	"github.com/kataras/iris/context"
)

func AuthService(Ctx context.Context)  {
	Ctx.Next()
}