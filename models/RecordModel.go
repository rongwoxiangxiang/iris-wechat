package models

import (
	"iris/common"
	"iris/config"
	"time"
)

type RecordModel struct {
	Id  int64 `xorm:"pk"`
	Wid int64
	Wuid int64
	Type string
	Content string
	CreatedAt time.Time `orm:"auto_now_add;type(datetime) created"`
}

func (this *RecordModel) TableName() string {
	return "records"
}

func (this *RecordModel) Insert() (insertId int64, err error){
	insertId, err = config.GetDb().InsertOne(this)
	if err != nil {
		err = common.ErrDataCreate
	}
	return
}

func (this *RecordModel) GetById() (RecordModel, error){
	if this.Id != 0{
		record := RecordModel{Id : this.Id}
		has, err := config.GetDb().Get(&record)
		if err != nil {
			err = common.ErrDataGet
		} else if has == false {
			err = common.ErrDataEmpty
		}
		return record,err
	}

	return RecordModel{}, common.ErrDataGet
}

//通过wid查找
func (r *RecordModel) LimitUnderWidList(index int,limit int) (records []RecordModel) {
	err := config.GetDb().Where("wid = ?",r.Wid).Limit(limit, (index - 1) * limit).Find(&records)
	if err != nil {
		err = common.ErrDataFind
	}
	return records
}

//通过wuid查找
func (r *RecordModel) LimitUnderWuidList(index int,limit int) (records []RecordModel) {
	err := config.GetDb().Where("wuid = ?",r.Wuid).Limit(limit, (index - 1) * limit).Find(&records)
	if err != nil {
		err = common.ErrDataFind
	}
	return records
}
