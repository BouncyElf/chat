package models

import (
	"strings"

	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

type Relation struct {
	UID string `gorm:"column:uid;primary_key" json:"uid"`

	// uids is the uid's friend list
	UIDs string `gorm:"column:uids" json:"uids"`
}

func (Relation) TableName() string {
	return "friend"
}

func (r *Relation) Save() {
	err := DB.Save(r).Error
	if err != nil {
		air.ERROR("save relation to db error", utils.M{
			"err":      err.Error(),
			"relation": r,
		})
	}
}

func NewRelation(uid string) *Relation {
	return &Relation{
		UID: uid,
	}
}

func GetRelation(uid string) *Relation {
	res := &Relation{}
	err := DB.Where("uid = ?", uid).Find(res).Error
	if err != nil {
		air.ERROR("get relation from db error", utils.M{
			"err": err.Error(),
			"uid": uid,
		})
		return nil
	}
	return res
}

func GetRelationsSlice(uids []string) []*Relation {
	res := []*Relation{}
	err := DB.Where("uid in (?)", uids).Find(&res).Error
	if err != nil {
		air.ERROR("get relations from db error", utils.M{
			"err":  err.Error(),
			"uids": uids,
		})
		return nil
	}
	return res
}

func GetRelations(uids []string) map[string]*Relation {
	res := map[string]*Relation{}
	relations := GetRelationsSlice(uids)
	if relations == nil {
		return res
	}
	for _, v := range relations {
		res[v.UID] = v
	}
	return res
}

func GetFriends(uid string) []string {
	relation := GetRelation(uid)
	if relation == nil {
		return []string{}
	}
	return strings.Split(relation.UIDs, ";")
}
