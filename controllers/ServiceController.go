package controllers

import (
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/chanxuehong/wechat/mp/message/callback/request"
	"github.com/chanxuehong/wechat/mp/message/callback/response"
	"github.com/kataras/iris"
	"iris/common"
	"iris/models"
	"log"
)

type ServiceController struct{
	Ctx iris.Context
}

var (
	msgHandler core.Handler
	msgServers  map[string]*core.Server
)

func init() {
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(defaultMsgHandler)
	mux.DefaultEventHandleFunc(defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, menuClickEventHandler)
	msgHandler = mux
	msgServers = make(map[string]*core.Server)
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
	msgServer := this.setMsgServer()
	msgServer.ServeHTTP(this.Ctx.ResponseWriter(), this.Ctx.Request(), nil)
}

func (this *ServiceController) PostIndex() {
	this.GetIndex()
}

func (this *ServiceController) setMsgServer() (msgServer *core.Server) {
	flag := this.Ctx.Params().Get("flag")
	if flag == "" {
		return
	}
	if server, ok := msgServers[flag]; ok == true {
		return server
	}
	wechat, err := (&models.WechatModel{Flag:flag}).GetByFlag()
	if err != nil {
		return
	}
	msgServer = core.NewServer("", wechat.Appid, wechat.Token, wechat.EncodingAesKey, msgHandler, nil)
	if wechat.NeedSaveMen == common.YES_VALUE {
		msgServers[flag] = msgServer
	}
	return
}