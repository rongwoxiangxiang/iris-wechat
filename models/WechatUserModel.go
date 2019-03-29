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
	CreatedAt time.Time `xorm:"created_at"`
	UpdatedAt time.Time `xorm:"updated_at"`
}

func (wu *WechatUserModel) TableName() string {
	return "wechat_users"
}

func (wu *WechatUserModel) Insert() (insertId int64, err error) {
	insertId, err = config.Orm.InsertOne(wu)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (wu *WechatUserModel) Update() (rows int64, err error){
	rows, err = config.Orm.Id(wu.Id).Update(wu)
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
	has, err := config.Orm.Where("wid= ? and openid = ?", wu.Wid, wu.Openid).Get(&user)
	if err != nil {
		return WechatUserModel{},common.ErrDataGet
	} else if has == false {
		return WechatUserModel{},common.ErrDataEmpty
	}
	return
}

func (wu *WechatUserModel) LimitUnderWidList(index int,limit int) (users []WechatUserModel) {
	err := config.Orm.Where("wid = ?",wu.Wid).Limit(limit, (index - 1) * limit).Find(&users)
	if err != nil {
		err = common.ErrDataFind
	}
	return users
}

func (wu *WechatUserModel) DeleteById() bool{
	_, err := config.Orm.Id(wu.Id).Unscoped().Delete(&WechatUserModel{})
	if err != nil{
		return false
	}
	return true
}