package models

import (
	"github.com/BouncyElf/chat/utils"
	"github.com/aofei/air"
)

type Unread struct {
	UID     string `gorm:"column:uid" json:"uid"`
	GID     string `gorm:"column:gid" json:"gid"`
	LastMID string `gorm:"column:last_mid" json:"last_mid"`
	Count   int    `gorm:"column:count" json:"count"`
}

func (Unread) TableName() string {
	return "unread"
}

func (u *Unread) Save() {
	err := DB.Save(u).Error
	if err != nil {
		air.ERROR("save unread to db error", utils.M{
			"err":    err.Error(),
			"unread": u,
		})
	}
}

func NewUnread(uid, gid string, lastMID ...string) *Unread {
	mid := ""
	if len(lastMID) != 0 {
		mid = lastMID[0]
	}
	return &Unread{
		UID:     uid,
		GID:     gid,
		LastMID: mid,
		Count:   0,
	}
}

func GetUnread(uid, gid string) *Unread {
	res := &Unread{}
	err := DB.Where("uid = ? AND gid = ?", uid, gid).Find(res).Error
	if err != nil {
		air.ERROR("get unread from db error", utils.M{
			"err": err.Error(),
			"uid": uid,
		})
		return nil
	}
	return res
}
