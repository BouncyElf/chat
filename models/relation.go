package models

import (
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

type Relation struct {
	UID  string `gorm:"column:uid;primary_key"`
	UIDs string `gorm:"column:uids"`
}

func (Relation) TableName() string {
	return "friend"
}

func (r *Relation) Save() {
	err := DB.Save(r).Error
	if err != nil {
		air.ERROR("save relation to db error", utils.M{
			"error":    err.Error(),
			"relation": r,
		})
	}
}

func NewRelation(uid string) *Relation {
	return &Relation{
		UID: uid,
	}
}

func GetRelations(uid string) *Relation {
	res := &Relation{}
	err := DB.Where("uid = ?", uid).Find(res).Error
	if err != nil {
		air.ERROR("get relation from db error", utils.M{
			"error": err.Error(),
			"uid":   uid,
		})
		return nil
	}
	return res
}
