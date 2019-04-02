package controllers

import (
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/chanxuehong/wechat/mp/message/callback/request"
	"github.com/chanxuehong/wechat/mp/message/callback/response"
	"github.com/chanxuehong/wechat/mp/user"
	"github.com/kataras/iris"
	"iris/common"
	"iris/models"
	"log"
	"time"
)

type ServiceController struct{
	Ctx iris.Context
}

type Service struct {
	*core.Server
	Wid int64
	AppSecret string
}

var (
	msgHandler core.Handler
	msgServers  map[string]*Service
)

func (this *ServiceController) init() {
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(this.defaultMsgHandler)
	mux.DefaultEventHandleFunc(this.defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, this.textMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, this.menuClickEventHandler)
	msgHandler = mux
	msgServers = make(map[string]*Service)
}

func (this *ServiceController) defaultMsgHandler(ctx *core.Context) {
	log.Printf("收到消息:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func (this *ServiceController) defaultEventHandler(ctx *core.Context) {
	log.Printf("收到事件:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func (this *ServiceController) textMsgHandler(ctx *core.Context) {
	log.Printf("收到文本消息:\n%s\n", ctx.MsgPlaintext)
	msg := request.GetText(ctx.MixedMsg)
	reply := models.ReplyModel{Wid: this.getWidFromMsgServers(), Alias: msg.Content}
	if reply.ActivityId == 0 {
		return
	}
	resp := this.doReplyTextAndClick(reply, msg.MsgHeader)
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func (this *ServiceController) menuClickEventHandler(ctx *core.Context) {
	log.Printf("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

	event := menu.GetClickEvent(ctx.MixedMsg)
	reply := models.ReplyModel{Wid: this.getWidFromMsgServers(), ClickKey: event.EventKey}
	if reply.ActivityId == 0 {
		return
	}
	resp := this.doReplyTextAndClick(reply, event.MsgHeader)
	//resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func (this *ServiceController) doReplyTextAndClick(reply models.ReplyModel, header core.MsgHeader) (msg interface{}) {
	wxUser := this.getWechatUser(header.FromUserName, reply.Wid)
	log.Print(wxUser)
	switch reply.Type {
	case models.REPLY_TYPE_TEXT:
		return response.NewText(header.FromUserName, header.ToUserName, header.CreateTime, reply.Success)
	case models.REPLY_TYPE_CODE:
		//return doReplyCode()
	case models.REPLY_TYPE_LUCKY:
		//return doReplyLucky()
	case models.REPLY_TYPE_CHECKIN:
		//return doReplyCheckin()
	}
	return 
}

//func (this *ServiceController) doReplyCode()  {
//
//}


func(this *ServiceController) getWechatUser(openId string, wid int64) (wechatUser models.WechatUserModel) {
	wechatUser.Wid = wid
	wechatUser.Openid = openId
	wechatUser,_ = wechatUser.GetByOpenid()
	go func(wechatUser models.WechatUserModel) {
		if wechatUser.Openid != "" && wechatUser.UpdatedAt.After(time.Now().Add(-24 * time.Hour)) {
			userInfo, err := user.Get(this.getWechatClient(), wechatUser.Openid, "")
			if err == nil {
				wechatUser.Nickname = userInfo.Nickname
				wechatUser.Sex = userInfo.Sex
				wechatUser.Province = userInfo.Province
				wechatUser.City = userInfo.City
				wechatUser.Country = userInfo.Country
				wechatUser.Language = userInfo.Language
				wechatUser.Headimgurl = userInfo.HeadImageURL
				wechatUser.Update()
			}
		}
	}(wechatUser)
	return
}

func (this *ServiceController) setMsgServer() (msgServer *core.Server) {
	flag := this.Ctx.Params().Get("flag")
	if flag == "" {
		return
	}
	if server, ok := msgServers[flag]; ok == true {
		return server.Server
	}
	wechat, err := (&models.WechatModel{Flag:flag}).GetByFlag()
	if err != nil {
		return
	}
	msgServer = core.NewServer("", wechat.Appid, wechat.Token, wechat.EncodingAesKey, msgHandler, nil)
	if wechat.NeedSaveMen == common.YES_VALUE {
		msgServers[flag] = &Service{Wid:wechat.Id, AppSecret:wechat.Appsecret, Server:msgServer}
	}
	return
}

func (this *ServiceController) getWechatClient() *core.Client {
	flag := this.Ctx.Params().Get("flag")
	accessTokenServer := core.NewDefaultAccessTokenServer(msgServers[flag].AppId(), msgServers[flag].AppSecret, nil)
	return core.NewClient(accessTokenServer, nil)
}

func (this *ServiceController) getWidFromMsgServers() int64 {
	flag := this.Ctx.Params().Get("flag")
	if flag == "" {
		return 0
	}
	return msgServers[flag].Wid
}

func (this *ServiceController) GetIndex() {
	msgServer := this.setMsgServer()
	msgServer.ServeHTTP(this.Ctx.ResponseWriter(), this.Ctx.Request(), nil)
}

func (this *ServiceController) PostIndex() {
	this.GetIndex()
}