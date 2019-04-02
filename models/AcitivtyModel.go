package models

import (
	"iris/common"
	"iris/config"
	"time"
)

type ActivityModel struct {
	Id  int64 `xorm:"pk"`
	Wid int64
	Name  string
	Desc  string
	Type int8
	events string
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
		return activity,err
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