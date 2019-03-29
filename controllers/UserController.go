package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"iris/middleware"
	"iris/models"
)

type UserController struct{
	Ctx iris.Context
}

func (c *UserController) GetAll() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}
func (c *UserController) GetPing() string {
	fmt.Print((&models.WechatUserModel{Id:1,Wid:4, Openid:"oSXCp1Omvs1-NVq5S2rsnpY_dPko"}).DeleteById())
	return "1111"
}

func (c *UserController) GetHello() interface{} {
	return map[string]string{"message": "Hello Iris!"}
}

func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/custom_path", "CustomHandlerWithoutFollowingTheNamingGuide", middleware.AuthService)
}

func (c *UserController) CustomHandlerWithoutFollowingTheNamingGuide() string {
	return "hello from the custom handler without following the naming guide"
}


func (c *UserController) Post() {}
func (c *UserController) Put() {}
func (c *UserController) Delete() {}
func (c *UserController) Connect() {}
func (c *UserController) Head() {}
func (c *UserController) Patch() {}
func (c *UserController) Options() {}
func (c *UserController) Trace() {}