package models

import (
	"github.com/BouncyElf/chat/utils"
	"github.com/aofei/air"
)

type List struct {
	UID  string `gorm:"column:uid;primary_key"`
	GIDs string `gorm:"column:gids"`
}

func (List) TableName() string {
	return "list"
}

func (l *List) Save() {
	err := DB.Save(l).Error
	if err != nil {
		air.ERROR("save list to db error", utils.M{
			"error": err.Error(),
			"list":  l,
		})
	}
}

func GetList(uid string) *List {
	res := &List{}
	err := DB.Where("uid = ?", uid).Find(res).Error
	if err != nil {
		air.ERROR("get list from db error", utils.M{
			"error": err.Error(),
			"uid":   uid,
		})
		return nil
	}
	return res
}

func NewList(uid string) *List {
	return &List{
		UID:  uid,
		GIDs: "",
	}
}
