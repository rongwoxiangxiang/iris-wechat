package models

import (
	"iris/common"
	"iris/config"
	"time"
)

const CHECK_FAIL = "签到失败，请重试！"

type CheckinModel struct {
	Id  int64 `xorm:"pk"`
	Wid int64
	ActivityId int64
	Wuid int64
	Liner int64
	Total int64
	Lastcheckin time.Time
	CreatedAt time.Time `orm:"auto_now_add;type(datetime) created"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime) updated"`
}

func (this *CheckinModel) TableName() string {
	return "checkins"
}

func (this *CheckinModel) Insert() (insertId int64, err error){
	insertId, err = config.GetDb().InsertOne(this)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (this *CheckinModel) DeleteById() bool{
	_, err := config.GetDb().Id(this.Id).Unscoped().Delete(&CheckinModel{})
	if err != nil{
		return false
	}
	return true
}

func (this *CheckinModel) List() (lotteries []CheckinModel) {
	if this.ActivityId == 0 || this.Wid == 0 {
		return
	}
	err := config.GetDb().Where("wid = ?", this.Wid).Find(&lotteries)
	if err != nil{
		lotteries = []CheckinModel{}
	}
	return lotteries
}

func (this *CheckinModel) GetCheckinByActivityWuid() (checkin CheckinModel, err error){
	if this.ActivityId == 0 || this.Wuid == 0 || this.Wid == 0{
		err = common.ErrDataGet
		return
	}
	checkin.ActivityId = this.ActivityId
	checkin.Wuid = this.Wuid
	has, err := config.GetDb().Get(&checkin)
	if err != nil {
		return CheckinModel{},common.ErrDataGet
	} else if has == false {
		checkin.ActivityId = this.ActivityId
		checkin.Wid = this.Wid
		checkin.Wuid = this.Wuid
		checkin.Lastcheckin = time.Now()
		checkin.Liner = 1
		checkin.Total = 1
		checkin.Id, err = checkin.Insert()
		if err != nil {
			return CheckinModel{},common.ErrDataCreate
		}
		return checkin,nil
	}
	return
}

func (this *CheckinModel) Update() (rows int64, err error){
	rows, err = config.GetDb().Id(this.Id).Update(this)
	if err != nil {
		err = common.ErrDataUpdate
	}
	return
}