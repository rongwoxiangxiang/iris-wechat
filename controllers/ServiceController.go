package controllers

import (
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/chanxuehong/wechat/mp/message/callback/request"
	"github.com/chanxuehong/wechat/mp/message/callback/response"
	"github.com/kataras/iris"
	"iris/config"
	"iris/models"
	"log"
)

type ServiceController struct{
	Ctx iris.Context
}


const (
	wxAppId     	= "wxac0647051de8c364"
	wxOriId         = "oriid"
	wxToken         = "chens"
	wxEncodedAESKey = "cH1XblyoXlx4Q9srC7HRWJaR0woP5Js2y1ujcwJTQa4"
)

var (
	msgHandler core.Handler
	msgServer  *core.Server
)

func init() {
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(defaultMsgHandler)
	mux.DefaultEventHandleFunc(defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, menuClickEventHandler)

	msgHandler = mux
}

func textMsgHandler(ctx *core.Context) {
	log.Printf("收到文本消息:\n%s\n", ctx.MsgPlaintext)
	msg := request.GetText(ctx.MixedMsg)
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultMsgHandler(ctx *core.Context) {
	log.Printf("收到消息:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func menuClickEventHandler(ctx *core.Context) {
	log.Printf("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultEventHandler(ctx *core.Context) {
	log.Printf("收到事件:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}


func (this *ServiceController) GetIndex() {
	this.setMsgServer()
	msgServer.ServeHTTP(this.Ctx.ResponseWriter(), this.Ctx.Request(), nil)
}

func (this *ServiceController) PostIndex() {
	this.GetIndex()
}

func (this *ServiceController) setMsgServer() {
	flag := this.Ctx.Params().Get("flag")
	if flag == "" {
		return
	}
	wechat := config.CacheGet(flag)
	if wechat == nil {
		wct, err := (&models.WechatModel{Flag:flag}).FindByFlag()
		if err != nil {
			return
		}
		config.CacheSet(flag, wct, 60)
	}

	msgServer = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)
}