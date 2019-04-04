package models

import (
	"iris/common"
	"iris/config"
	"time"
)

type PrizeHistoryModel struct {
	Id  int64 `xorm:"pk"`
	ActivityId int64
	Wuid int64
	Prize string
	Code string
	Level int8
	CreatedAt time.Time `orm:"auto_now_add;type(datetime) created"`
}

func (this *PrizeHistoryModel) TableName() string {
	return "prize_history"
}

func (this *PrizeHistoryModel) Insert() (insertId int64, err error){
	insertId, err = config.GetDb().InsertOne(this)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (this *PrizeHistoryModel) GetByActivityWuId() (PrizeHistoryModel, error){
	history := PrizeHistoryModel{ActivityId : this.ActivityId, Wuid:this.Wuid}
	has, err := config.GetDb().Desc("id").Get(&history)
	if err != nil {
		err = common.ErrDataGet
	} else if has == false {
		err = common.ErrDataEmpty
	}
	return history,err
}

func (this *PrizeHistoryModel) DeleteById() bool{
	_, err := config.GetDb().Id(this.Id).Unscoped().Delete(&PrizeHistoryModel{})
	if err != nil{
		return false
	}
	return true
}

func (this *PrizeHistoryModel) LimitUnderActivityList(index int,limit int) (histories []PrizeHistoryModel) {
	err := config.GetDb().Where("acitivity_id = ?", this.ActivityId).Limit(limit, (index - 1) * limit).Find(&histories)
	if err != nil {
		err = common.ErrDataFind
	}
	return histories
}
