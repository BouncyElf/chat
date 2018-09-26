package handlers

import (
	"errors"
	"strings"

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
}

func newGroupHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	tuids, _ := req.Params["tuids"]
	groupType, _ := req.Params["group_type"]
	groupName, _ := req.Params["group_name"]
	go models.NewGroup(append(strings.Split(tuids, ";"), uid),
		groupType, groupName).Save()
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
	for _, v := range strings.Split(group.UIDs, ";") {
		if v == uid {
			return true
		}
	}
	return false
}
