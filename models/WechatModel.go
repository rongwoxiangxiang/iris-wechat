package models

import (
	"iris/common"
	"iris/config"
	"time"
)

type WechatModel struct {
	Id int64 `xorm:"pk"`
	Gid int64
	Name string
	Appid string
	Appsecret string
	EncodingAesKey string
	Token string
	Flag  string
	Type  int8
	Pass  int8
	SaveInput int8
	NeedSaveMen int8	//该公众号hander是否持久化
	CreatedAt time.Time `orm:"auto_now_add;type(datetime) created"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime) updated"`
}

func (w *WechatModel) TableName() string {
	return "wechats"
}

func (w *WechatModel) Insert() (insertId int64, err error){
	insertId, err = config.GetDb().InsertOne(w)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (w *WechatModel) GetById() (WechatModel, error){
	if w.Id != 0{
		wechat := WechatModel{Id : w.Id}
		has, err := config.GetDb().Get(&wechat)
		if err != nil {
			err = common.ErrDataGet
		} else if has == false {
			err = common.ErrDataEmpty
		}
		return wechat,err
	}
	return WechatModel{},common.ErrDataGet
}

func (w *WechatModel) DeleteById() bool{
	_, err := config.GetDb().Id(w.Id).Unscoped().Delete(&WechatModel{})
	if err != nil{
		return false
	}
	return true
}

func (w *WechatModel) FindByGid() (wechats []WechatModel) {
	err := config.GetDb().Where("gid = ?",w.Gid).Find(&wechats)
	if err != nil {
		err = common.ErrDataFind
	}
	return wechats
}

func (w *WechatModel) GetByAppid() (WechatModel, error) {
	wechat := &WechatModel{Appid:w.Appid}
	has, err := config.GetDb().Get(wechat)
	if err != nil {
		err = common.ErrDataGet
	} else if has == false {
		err = common.ErrDataEmpty
	}
	return *wechat,err
}

func (w *WechatModel) GetByFlag() (WechatModel, error) {
	wechat := &WechatModel{}
	has, err := config.GetDb().Where("flag = ?",w.Flag).Get(wechat)
	if err != nil {
		err = common.ErrDataGet
	} else if has == false {
		err = common.ErrDataEmpty
	}
	return *wechat,err
}
