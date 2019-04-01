package models

import (
	"iris/common"
	"iris/config"
	"time"
)

const (
	REPLY_TYPE_TEXT = "text"
	REPLY_TYPE_CODE = "code"
	REPLY_TYPE_LUCKY = "luck"
	REPLY_TYPE_CHECKIN = "checkin"
)

const PLEASE_TRY_AGAIN = "活动太火爆了，请稍后重试"

type ReplyModel struct {
	Id  int64
	Wid int64
	ActivityId int64
	Alias string
	ClickKey string
	Success string
	Fail string
	Extra string
	Type string
	Disabled int8
	Match int8
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *ReplyModel) TableName() string {
	return "replies"
}

func (w *ReplyModel) Insert() (insertId int64, err error){
	insertId, err = config.GetDb().InsertOne(w)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (r *ReplyModel) GetById() (ReplyModel, error){
	if r.Id != 0{
		user := ReplyModel{Id : r.Id}
		has, err := config.GetDb().Get(&user)
		if err != nil {
			err = common.ErrDataGet
		} else if has == false {
			err = common.ErrDataEmpty
		}
		return user,err
	}

	return ReplyModel{}, common.ErrDataGet
}

func (r *ReplyModel) DeleteById() bool{
	_, err := config.GetDb().Id(r.Id).Unscoped().Delete(&ReplyModel{})
	if err != nil{
		return false
	}
	return true
}

/**
 * @Find
 * @Param Reply.Id int
 * @Param Reply.Alias string
 * @Param Reply.ClickKey string
 * @Success Reply
 */
func (r *ReplyModel) FindOne() (reply ReplyModel) {
	if "" == r.Alias && "" == r.ClickKey {
		return
	}
	qs := config.GetDb().Where("wid = ?",r.Wid)
	if "" != r.Alias {
		qs = qs.Where("alias = ?", r.Alias)
	} else if "" != r.ClickKey {
		qs = qs.Where("click_key = ?", r.ClickKey)
	}
	_, err := qs.Where("disabled = ?", common.NO_VALUE).Get(&reply)
	if err != nil{
		return ReplyModel{}
	}
	return reply
}

func (r *ReplyModel) LimitUnderWidList(index int,limit int) (relpies []ReplyModel) {
	err := config.GetDb().Where("wid = ?",r.Wid).Limit(limit, (index - 1) * limit).Find(&relpies)
	if err != nil {
		err = common.ErrDataFind
	}
	return relpies
}
