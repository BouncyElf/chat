package models

import (
	"strings"

	"github.com/BouncyElf/chat/common"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

type Group struct {
	// snowflake
	GID string `gorm:"column:gid;primary_key" json:"gid"`

	// private chat, name is `a;b`
	Name string `gorm:"colum:name" json:"name"`
	// uid1;uid2;uid3
	UIDs string `gorm:"column:uids" json:"uids"`
	Type string `gorm:"column:type" json:"type"`
}

func (Group) TableName() string {
	return "group"
}

func (g *Group) Save() {
	if g.GID == "" {
		g.GID = common.NewSnowFlake()
	}
	err := DB.Save(g).Error
	if err != nil {
		air.ERROR("save group to db error", utils.M{
			"err":   err.Error(),
			"group": g,
		})
	}
}

func (g *Group) Delete() {
	if g.GID == "" {
		air.ERROR("attempt to delete group without id", utils.M{
			"group": g,
		})
		return
	}
	err := DB.Delete(g).Error
	if err != nil {
		air.ERROR("delete group error", utils.M{
			"err":   err.Error(),
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

func GetGroup(gid string) *Group {
	res := &Group{}
	err := DB.Where("gid = ?", gid).Find(res).Error
	if err != nil {
		air.ERROR("get group from db error", utils.M{
			"err": err.Error(),
			"gid": gid,
		})
		return nil
	}
	return res
}

func GetGroups(gids []string) map[string]*Group {
	groups := GetGroupsSlice(gids)
	res := map[string]*Group{}
	for _, v := range groups {
		res[v.GID] = v
	}
	return res
}

func GetGroupsSlice(gids []string) []*Group {
	groups := []*Group{}
	err := DB.Where("gid in (?)", gids).Find(groups).Error
	if err != nil {
		air.ERROR("get groups error", utils.M{
			"err":  err.Error(),
			"gids": gids,
		})
	}
	return groups
}

func GetPrivateChatGroup(uid, tuid string) *Group {
	uids := strings.Join([]string{uid, tuid}, ";")
	group := &Group{}
	err := DB.Where("uids = ? AND type = ?", uids, common.ChatTypePrivate).
		Find(group).Error
	if err != nil {
		air.ERROR("get group from db error", utils.M{
			"err":  err.Error(),
			"uid":  uid,
			"tuid": tuid,
			"uids": uids,
		})
		return nil
	}
	if group.GID == "" {
		uids = strings.Join([]string{tuid, uid}, ";")
		err = DB.Where("uids = ? AND type = ?", uids,
			common.ChatTypePrivate).Find(group).Error
		if err != nil {
			air.ERROR("get group from db error", utils.M{
				"err":  err.Error(),
				"uid":  uid,
				"tuid": tuid,
				"uids": uids,
			})
			return nil
		}
	}
	return group
}
