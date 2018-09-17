package models

import (
	"strings"

	"github.com/BouncyElf/chat/common"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

type Group struct {
	// snowflake
	GID int64 `gorm:"column:gid;primary_key"`

	// private chat, name is `a;b`
	Name string `gorm:"colum:name"`
	// uid1;uid2;uid3
	UIDs string `gorm:"column:uids"`
	Type string `gorm:"column:type"`
}

func (Group) TableName() string {
	return "group"
}

func (g *Group) Save() {
	if g.GID == 0 {
		g.GID = common.NewSnowFlake()
	}
	err := DB.Save(g).Error
	if err != nil {
		air.ERROR("save group to db error", utils.M{
			"error": err.Error(),
			"group": g,
		})
	}
}

func NewGroup(uids []string, name, t string) *Group {
	return &Group{
		GID:  common.NewSnowFlake(),
		Name: name,
		UIDs: strings.Join(uids, ";"),
		Type: t,
	}
}

func GetGroup(gid int64) *Group {
	res := &Group{}
	err := DB.Where("gid = ?", gid).Find(res).Error
	if err != nil {
		air.ERROR("get group from db error", utils.M{
			"error": err.Error(),
			"gid":   gid,
		})
		return nil
	}
	return res
}
