package handlers

import (
	"errors"
	"strings"

	"github.com/BouncyElf/chat/common"
	"github.com/BouncyElf/chat/gas"
	"github.com/BouncyElf/chat/models"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

func init() {
	a := air.Group{
		Prefix: "/api/group",
		Gases:  []air.Gas{gas.Auth},
	}
	a.POST("/new", newGroupHandler)
	a.POST("/addmem", addMemberHandler)
	a.POST("/update/name", updateGroupNameHandler)
}

func newGroupHandler(req *air.Request, res *air.Response) error {
	tuids, _ := req.Params["tuids"]
	groupName, _ := req.Params["group_name"]
	if groupName == "" {
		groupName = common.DefaultGroupName
	}
	group := models.NewGroup(
		strings.Split(tuids, ";"),
		common.ChatTypeGroup,
		groupName,
	)
	group.Save()
	for _, uid := range strings.Split(tuids, ";") {
		if !inList(uid, group.GID) {
			UpdateListAdd(uid, group.GID)
		}
	}
	return utils.Success(res, "")
}

func addMemberHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	tuid := req.Params["tuid"]
	gid := req.Params["gid"]
	group := models.GetGroup(gid)
	if group == nil {
		air.ERROR("group not found", utils.M{
			"uid":  uid,
			"tuid": tuid,
			"gid":  gid,
		})
		return utils.Error(404, errors.New("group not found"))
	}
	if !strings.Contains(group.UIDs, tuid) {
		group.UIDs = strings.Join([]string{group.UIDs, tuid}, ";")
		go group.Save()
	}
	return utils.Success(res, "")
}

func exitGroupHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	gid := req.Params["gid"]
	group := models.GetGroup(gid)
	if group == nil {
		air.ERROR("group not found", utils.M{
			"uid": uid,
			"gid": gid,
		})
		return utils.Error(404, errors.New("group not found"))
	}
	if strings.Contains(group.UIDs, uid) {
		uids := strings.Split(group.UIDs, ";")
		for k, v := range uids {
			if v == uid {
				uids = append(uids[:k], uids[k+1:]...)
				break
			}
		}
		group.UIDs = strings.Join(uids, ";")
		go group.Save()
	}
	return utils.Success(res, "")
}

func IsInGroup(uid string, gid string) bool {
	group := models.GetGroup(gid)
	if group == nil {
		return false
	}
	return strings.Contains(group.UIDs, uid)
}

func updateGroupNameHandler(req *air.Request, res *air.Response) error {
	gid := req.Params["gid"]
	uid := req.Params["uid"]
	name := req.Params["name"]
	if !IsInGroup(uid, gid) {
		air.ERROR("not in group", utils.M{
			"gid":  gid,
			"uid":  uid,
			"name": name,
		})
		return utils.Error(400, errors.New("bad request"))
	}
	group := models.GetGroup(gid)
	if group == nil {
		air.ERROR("group not found", utils.M{
			"gid":  gid,
			"uid":  uid,
			"name": name,
		})
		return utils.Error(404, errors.New("group not found"))
	}
	group.Name = name
	go group.Save()
	return utils.Success(res, "")
}
