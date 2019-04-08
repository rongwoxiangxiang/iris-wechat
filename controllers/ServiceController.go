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
	"iris/modules/ai"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	wxflag string
	msgHandler core.Handler
	msgServers  map[string]*Service
)

func init() {
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(defaultMsgHandler)
	mux.DefaultEventHandleFunc(defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler, defaultTextMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, menuClickEventHandler)
	msgHandler = mux
	msgServers = make(map[string]*Service)
}

func defaultMsgHandler(ctx *core.Context) {
	log.Printf("收到消息:\n%s\n", )
	ctx.NoneResponse()
}

func defaultEventHandler(ctx *core.Context) {
	log.Printf("收到事件:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func defaultTextMsgHandler(ctx *core.Context) {
	log.Printf("AI智能闲聊:\n%s\n", ctx.MsgPlaintext)
	msg := request.GetText(ctx.MixedMsg)
	aiService := ai.Ai{}
	aiService.SetAiServer()
	answer := aiService.NlpTextchat(&http.Client{}, msg.Content, msg.FromUserName)
	if answer.AnswerData != "" {
		resp := response.NewText(msg.MsgHeader.FromUserName, msg.MsgHeader.ToUserName, msg.MsgHeader.CreateTime, answer.AnswerData)
		if len(ctx.AESKey) == 0 {
			ctx.RawResponse(resp)
			return
		}
		ctx.AESResponse(resp, 0, "", nil)
		return
	}
	ctx.NoneResponse()
}

func textMsgHandler(ctx *core.Context) {
	log.Printf("收到文本消息:\n%s\n", ctx.MsgPlaintext)
	msg := request.GetText(ctx.MixedMsg)
	wxUser := getWechatUser(msg.FromUserName, ctx.MixedMsg)
	reply := (&models.ReplyModel{Wid: wxUser.Wid, Alias: msg.Content}).FindOne()
	if reply.Success == "" {
		return
	}
	resp := responseTextAndClick(reply, wxUser, msg.MsgHeader)
	if len(ctx.AESKey) == 0 {
		ctx.RawResponse(resp)
		return
	}
	ctx.AESResponse(resp, 0, "", nil)
}

func menuClickEventHandler(ctx *core.Context) {
	log.Printf("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)
	event := menu.GetClickEvent(ctx.MixedMsg)
	wxUser := getWechatUser(event.FromUserName, ctx.MixedMsg)
	reply := (&models.ReplyModel{Wid: wxUser.Wid, ClickKey: event.EventKey}).FindOne()
	if reply.ActivityId == 0 {
		return
	}
	resp := responseTextAndClick(reply, wxUser, event.MsgHeader)
	if len(ctx.AESKey) == 0 {
		ctx.RawResponse(resp)
		return
	}
	ctx.AESResponse(resp, 0, "", nil)
}

func responseTextAndClick(reply models.ReplyModel, wxUser models.WechatUserModel, header core.MsgHeader) (msg interface{}) {
	switch reply.Type {
	case models.REPLY_TYPE_TEXT:
		return response.NewText(
			header.FromUserName,
			header.ToUserName,
			header.CreateTime,
			reply.Success)
	case models.REPLY_TYPE_CODE:
		return response.NewText(
			header.FromUserName,
			header.ToUserName,
			header.CreateTime,
			doReplyCode(reply, wxUser))
	case models.REPLY_TYPE_CHECKIN:
		return response.NewText(
			header.FromUserName,
			header.ToUserName,
			header.CreateTime,
			doReplyCheckin(reply, wxUser))
	case models.REPLY_TYPE_LUCKY:
		return response.NewText(
			header.FromUserName,
			header.ToUserName,
			header.CreateTime,
			doReplyLucky(reply, wxUser))
	}
	return 
}

func doReplyCode(reply models.ReplyModel, wxUser models.WechatUserModel) string {
	history, _ := (&models.PrizeHistoryModel{ActivityId:reply.ActivityId, Wuid:wxUser.Id}).GetByActivityWuId()
	if history.Code != "" {
		return strings.Replace(reply.Success, "%code%", history.Code, 1)
	}
	prize, err := (&models.PrizeModel{ActivityId:reply.ActivityId, Level:int8(models.PRIZE_LEVEL_DEFAULT)}).GetOneUsedPrize(0)
	if err == common.ErrDataUnExist {
		return reply.NoPrizeReturn
	}
	if prize.Code != "" {
		_, err = (&models.PrizeHistoryModel{ActivityId:reply.ActivityId, Wuid:wxUser.Id, Code:prize.Code}).Insert()
		if err != nil {
			return models.PLEASE_TRY_AGAIN
		}
		return strings.Replace(reply.Success, "%code%", prize.Code, 1)
	}
	return models.PLEASE_TRY_AGAIN
}

func doReplyCheckin(reply models.ReplyModel, wxUser models.WechatUserModel) string 	{
	checkin, err := (&models.CheckinModel{ActivityId:reply.ActivityId, Wuid:wxUser.Id, Wid:wxUser.Wid}).GetCheckinByActivityWuid()
	if err != nil {
		return models.CHECK_FAIL
	}
	lastCheckinDate := checkin.Lastcheckin.Format("2006-01-02")
	if lastCheckinDate == time.Now().Format("2006-01-02") {
		return strings.
			NewReplacer("%liner%",  strconv.FormatInt(checkin.Liner, 10), "%total%", strconv.FormatInt(checkin.Total, 10)).
			Replace(reply.Success)
	}
	if lastCheckinDate == time.Now().Add(-24 * time.Hour).Format("2006-01-02"){//连续签到
		checkin.Liner = checkin.Liner + 1
	} else {//重置连续签到数
		checkin.Liner = 1
	}
	checkin.Total = checkin.Total + 1
	checkin.Lastcheckin = time.Now()
	_, err = checkin.Update()
	if err != nil {
		return models.CHECK_FAIL
	}
	return strings.
		NewReplacer("%liner%", strconv.FormatInt(checkin.Liner, 10), "%total%", strconv.FormatInt(checkin.Total, 10)).
		Replace(reply.Success)
}

func doReplyLucky(reply models.ReplyModel, wxUser models.WechatUserModel) string {
	activity, err := (&models.ActivityModel{Id:reply.ActivityId}).GetById()
	now := time.Now()
	if err != nil || activity.TimeStarted.IsZero() || activity.TimeEnd.IsZero() {
		return models.ACTIVITY_DATA_UNDEFINE
	} else if now.Before(activity.TimeStarted){
		return models.ACTIVITY_DATE_BEFORE
	} else if now.After(activity.TimeEnd) {
		return models.ACTIVITY_DATE_AFTER
	}
	history, _ := (&models.PrizeHistoryModel{ActivityId: reply.ActivityId, Wuid: wxUser.Id}).GetByActivityWuId()
	if activity.ActivityType == models.ACTIVITY_TYPE_LUCK_CHECKIN {//签到抽奖，验证签到条件是否满足
		if history.Prize != "" {
			return strings.NewReplacer("%prize%",  history.Prize, "%code%", history.Code).Replace(reply.Success)
		}
		checkin, err := (&models.CheckinModel{ActivityId:activity.RelativeId,Wuid:wxUser.Id}).GetCheckinUserData()
		if err != nil {
			return reply.Fail
		}
		if need, _ := strconv.ParseInt(activity.Extra, 10 , 64); checkin.Total < need {
			return reply.Fail
		}
	} else if activity.ActivityType == models.ACTIVITY_TYPE_LUCK_EVERYDAY {//每日抽奖，验证今日是否已经获取
		if history.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02") {
			return strings.NewReplacer("%prize%",  history.Prize, "%code%", history.Code).Replace(reply.Success)
		}
	} else {
		if history.Prize != "" {
			return strings.NewReplacer("%prize%",  history.Prize, "%code%", history.Code).Replace(reply.Success)
		}
	}

	luck, err := (&models.LotteryModel{Wid:reply.Wid, ActivityId:reply.ActivityId}).Luck()
	if err == common.ErrLuckFinal {
		return common.ErrLuckFinal.Msg
	} else if err == common.ErrDataUnExist {
		return reply.NoPrizeReturn
	} else if err != nil {
		return common.ErrLuckFail.Msg
	}
	var prize models.PrizeModel
	if luck.FirstCodeId != 0 {//表示存在礼包码
		prize, err = (&models.PrizeModel{ActivityId:reply.ActivityId, Level:luck.Level}).GetOneUsedPrize(luck.FirstCodeId)
		if err == common.ErrDataUnExist {
			return models.PLEASE_TRY_AGAIN
		}
	}
	if luck.Name != "" {
		(&models.PrizeHistoryModel{ActivityId:reply.ActivityId, Wuid:wxUser.Id, Prize:luck.Name, Code:prize.Code, Level:luck.Level}).Insert()
	}
	return strings.NewReplacer("%prize%",  luck.Name, "%code%", prize.Code).Replace(reply.Success)
}

func getWechatUser(openId string, msg *core.MixedMsg) (wechatUser models.WechatUserModel) {
	wechatUser.Wid = msgServers[wxflag].Wid
	wechatUser.Openid = openId
	wechatUser,_ = wechatUser.GetByOpenid()
	go func(wechatUser models.WechatUserModel, msg *core.MixedMsg) {
		(&models.RecordModel{
			Wid:wechatUser.Wid,
			Wuid:wechatUser.Id,
			Type: string(msg.MsgType) + string(msg.EventType),
			Content: msg.Content + msg.EventKey + msg.MediaId,
		}).Insert()//存储用户操作
		if wechatUser.Openid != "" && (wechatUser.UpdatedAt.IsZero() || time.Now().After(wechatUser.UpdatedAt.Add(24 * time.Hour))) {
			userInfo, err := user.Get(getWechatClient(wxflag), wechatUser.Openid, "")
			if err == nil {
				(&models.WechatUserModel{
					Id:wechatUser.Id,
					Nickname : userInfo.Nickname,
					Sex : userInfo.Sex,
					Province : userInfo.Province,
					City : userInfo.City,
					Country : userInfo.Country,
					Language : userInfo.Language,
					Headimgurl : userInfo.HeadImageURL,
					UpdatedAt : time.Now(),
				}).Update()
			}
		}
	}(wechatUser, msg)
	return
}

func setMsgServer(flag string) (msgServer *core.Server) {
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

func getWechatClient(flag string) *core.Client {
	accessTokenServer := core.NewDefaultAccessTokenServer(msgServers[flag].AppId(), msgServers[flag].AppSecret, nil)
	return core.NewClient(accessTokenServer, nil)
}

func (this *ServiceController) GetIndex() {
	flag := this.Ctx.Params().Get("flag")
	if flag == "" {
		return
	}
	wxflag = flag
	msgServer := setMsgServer(flag)
	msgServer.ServeHTTP(this.Ctx.ResponseWriter(), this.Ctx.Request(), nil)
}

func (this *ServiceController) PostIndex() {
	this.GetIndex()
}