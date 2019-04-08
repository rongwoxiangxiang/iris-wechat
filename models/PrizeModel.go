package models

import (
	"iris/common"
	"iris/config"
	"time"
)

const (
	PRIZE_LEVEL_DEFAULT = 0
)

type PrizeModel struct {
	Id  int64 `xorm:"pk"`
	Wid int64
	ActivityId int64
	Code string
	Level int8
	Used int8
	CreatedAt time.Time `orm:"auto_now_add;type(datetime) created"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime) updated"`
}

func (this *PrizeModel) TableName() string {
	return "prizes"
}

func (this *PrizeModel) GetOneUsedPrize(idGt int64) (prize PrizeModel,err error) {
	qs := config.GetDb()
	if idGt > 0 {
		qs.Where("id >= ?", idGt)
	}
	has, err := qs.Where("activity_id = ? AND level = ? AND used = ?", this.ActivityId, this.Level, common.NO_VALUE).Get(&prize)
	if err != nil || has == false {
		err = common.ErrDataUnExist
		return
	}
	prize.Used = common.YES_VALUE
	_, err = config.GetDb().
		Where("id = ? and used = ?", prize.Id, common.NO_VALUE).
		Cols("used").
		Update(&prize)
	if err != nil {
		err = common.ErrDataUpdate
		return
	}
	return
}

func (this *PrizeModel) Insert() (insertId int64, err error){
	insertId, err = config.GetDb().InsertOne(this)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (this *PrizeModel) InsertBatch(prizes []PrizeModel) (error){
	_, err := config.GetDb().Insert(&prizes)
	if err != nil{
		return common.ErrDataCreate
	}
	return nil
}

func (this *PrizeModel) DeleteById() bool{
	_, err := config.GetDb().Id(this.Id).Unscoped().Delete(&PrizeModel{})
	if err != nil{
		return false
	}
	return true
}

func (this *PrizeModel) LimitUnderActivityList(index int,limit int) (prizes []PrizeModel) {
	err := config.GetDb().Where("acitivity_id = ?", this.ActivityId).Limit(limit, (index - 1) * limit).Find(&prizes)
	if err != nil {
		err = common.ErrDataFind
	}
	return prizes
}
