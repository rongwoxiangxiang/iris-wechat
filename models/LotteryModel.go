package models

import (
	"iris/common"
	"iris/config"
	"math/rand"
	"time"
)

const MAX_LUCKY_NUM = 10000

type LotteryModel struct {
	Id  int64 `xorm:"pk"`
	Wid int64
	ActivityId int64
	Name string
	Desc string
	TotalNum int64
	ClaimedNum int64
	Probability int
	FirstCodeId int64
	Level int8
	CreatedAt time.Time `orm:"auto_now_add;type(datetime) created"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime) updated"`
}

func (this *LotteryModel) TableName() string {
	return "lotteries"
}

func (this *LotteryModel) Insert() (insertId int64, err error){
	insertId, err = config.GetDb().InsertOne(this)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (this *LotteryModel) DeleteById() bool{
	_, err := config.GetDb().Id(this.Id).Unscoped().Delete(&LotteryModel{})
	if err != nil{
		return false
	}
	return true
}

func (this *LotteryModel) List() (lotteries []LotteryModel) {
	if this.ActivityId == 0 || this.Wid == 0 {
		return
	}
	err := config.GetDb().Where("wid = ? and activity_id = ?", this.Wid, this.ActivityId).Find(&lotteries)
	if err != nil{
		lotteries = []LotteryModel{}
	}
	return lotteries
}

//抽奖
func (this *LotteryModel) Luck() (lottery LotteryModel, err error) {
	lotteries := this.List()
	if len(lotteries) < 1 {
		err = common.ErrDataUnExist
		return
	}
	max := MAX_LUCKY_NUM
	actvityFinished := true
	for _, lot := range lotteries {
		if lot.ClaimedNum >= lot.TotalNum {//当前奖品发放完毕
			max -= lot.Probability
			continue
		}
		actvityFinished = false
		random := rand.Intn(max)
		if random <= lot.Probability {
			lottery = lot
			break
		}
		max -= lot.Probability
	}
	if actvityFinished {//全部奖品发放完毕，自动结束活动
		(&ReplyModel{Wid: this.Wid, ActivityId: this.ActivityId}).ChangeDisabledByWidActivityId(common.YES_VALUE)
		err = common.ErrLuckFinal
		return
	}

	if lottery.Id == 0 {
		err = common.ErrLuckFail
		return
	}
	claimedNum := lottery.ClaimedNum
	lottery.ClaimedNum = claimedNum + 1
	_, err = config.GetDb().
		Table(new(LotteryModel)).
		Where("id = ? and claimed_num = ?", lottery.Id, claimedNum).
		Cols("claimed_num").
		Update(lottery)
	if err != nil {
		err = common.ErrDataUpdate
		return
	}
	return
}