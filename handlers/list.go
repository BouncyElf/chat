package handlers

import (
	"strings"

	"github.com/BouncyElf/chat/gas"
	"github.com/BouncyElf/chat/models"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

func init() {
	a := air.Group{
		Prefix: "/api/list",
		Gases:  []air.Gas{gas.Auth},
	}
	a.POST("/get", getListHandler)
}

func getListHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	list := models.GetList(uid)
	if list == nil {
		models.NewList(uid).Save()
		return utils.Success(res, "")
	}
	groups := models.GetGroupsSlice(strings.Split(list.GIDs, ";"))
	data := []interface{}{}
	for _, group := range groups {
		data = append(data, utils.M{
			"gid":    group.GID,
			"name":   group.Name,
			"uids":   group.UIDs,
			"type":   group.Type,
			"unread": hasUnreadMsg(uid, group.GID),
		})
	}
	return utils.Success(res, data)
}

func inList(uid, gid string) bool {
	list := models.GetList(uid)
	if list == nil {
		return false
	}
	return strings.Contains(list.GIDs, gid)
}

func UpdateListAdd(uid, gid string) {
	list := models.GetList(uid)
	if list == nil {
		list = models.NewList(uid)
	}
	if !strings.Contains(list.GIDs, gid) {
		list.GIDs = strings.Join([]string{list.GIDs, gid}, ";")
	}
	go list.Save()
}
