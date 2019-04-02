package models

import (
	"iris/common"
	"iris/config"
	"time"
)

type WechatUserModel struct {
	Id  int64 `xorm:"pk Int"`
	Wid int64 `xorm:"wid"`
	UserId  int64 `xorm:"user_id"`
	Openid  string `xorm:"varchar(64)"`
	Nickname string `xorm:"varchar(64)"`
	Sex int
	Province string `orm:"varchar(20)"`
	City string `orm:"varchar(20)"`
	Country string `orm:"varchar(20)"`
	Language string `orm:"varchar(20)"`
	Headimgurl string `orm:"varchar(200)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime) created"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime) updated"`
}

func (wu *WechatUserModel) TableName() string {
	return "wechat_users"
}

func (wu *WechatUserModel) Insert() (insertId int64, err error) {
	insertId, err = config.GetDb().InsertOne(wu)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (wu *WechatUserModel) Update() (rows int64, err error){
	rows, err = config.GetDb().Id(wu.Id).Update(wu)
	if err != nil {
		err = common.ErrDataUpdate
	}
	return
}

func (wu *WechatUserModel) GetByOpenid() (user WechatUserModel, err error){
	if wu.Openid == "" || wu.Wid == 0{
		err = common.ErrDataGet
		return
	}
	user.Wid = wu.Wid
	user.Openid = wu.Openid
	has, err := config.GetDb().Get(&user)
	if err != nil {
		return WechatUserModel{},common.ErrDataGet
	} else if has == false {
		user.Id, err = user.Insert()
		if err != nil {
			return WechatUserModel{},common.ErrDataCreate
		}
		return user,nil
	}
	return
}

func (wu *WechatUserModel) LimitUnderWidList(index int,limit int) (users []WechatUserModel) {
	err := config.GetDb().Where("wid = ?",wu.Wid).Limit(limit, (index - 1) * limit).Find(&users)
	if err != nil {
		err = common.ErrDataFind
	}
	return users
}

func (wu *WechatUserModel) DeleteById() bool{
	_, err := config.GetDb().Id(wu.Id).Unscoped().Delete(&WechatUserModel{})
	if err != nil{
		return false
	}
	return true
}