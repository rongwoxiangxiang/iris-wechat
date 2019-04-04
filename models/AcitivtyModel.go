package models

import (
	"iris/common"
	"iris/config"
	"time"
)

const (
	ACTIVITY_DATA_UNDEFINE = "活动信息不存在！"
	ACTIVITY_DATE_BEFORE   = "活动未开始！"
	ACTIVITY_DATE_AFTER    = "活动已结束！"
)

const (
	ACTIVITY_TYPE_LUCK_DIRECT   = 11 //直接抽奖,单次
	ACTIVITY_TYPE_LUCK_CHECKIN  = 12 //签到抽奖，签到固定天数后抽奖
	ACTIVITY_TYPE_LUCK_EVERYDAY = 13 //每日抽奖，每天一次

	ACTIVITY_TYPE_CODE 			= 21 //直接发放奖励
	ACTIVITY_TYPE_CHECKIN 		= 31 //签到
)

type ActivityModel struct {
	Id  int64 `xorm:"pk"`
	Wid int64
	Name  string
	Desc  string
	ActivityType int8
	RelativeId int64
	Extra string
	TimeStarted time.Time
	TimeEnd time.Time
	CreatedAt time.Time `orm:"auto_now_add;type(datetime) created"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime) updated"`
}

func (this *ActivityModel) TableName() string {
	return "activities"
}

func (this *ActivityModel) GetById() (ActivityModel, error){
	if this.Id != 0{
		activity := ActivityModel{Id : this.Id}
		has, err := config.GetDb().Get(&activity)
		if err != nil {
			err = common.ErrDataGet
		} else if has == false {
			err = common.ErrDataEmpty
		}
		return activity, err
	}
	return ActivityModel{},common.ErrDataGet
}

func (this *ActivityModel) Insert() (insertId int64, err error){
	insertId, err = config.GetDb().InsertOne(this)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (this *ActivityModel) DeleteById() bool{
	_, err := config.GetDb().Id(this.Id).Unscoped().Delete(&ActivityModel{})
	if err != nil{
		return false
	}
	return true
}

func (this *ActivityModel) LimitUnderWidList(index int,limit int) (activities []ActivityModel) {
	err := config.GetDb().Where("wid = ?", this.Wid).Limit(limit, (index - 1) * limit).Find(&activities)
	if err != nil {
		err = common.ErrDataFind
	}
	return activities
}

func (this *ActivityModel) Update() (rows int64, err error){
	rows, err = config.GetDb().Id(this.Id).Update(this)
	if err != nil {
		err = common.ErrDataUpdate
	}
	return
}