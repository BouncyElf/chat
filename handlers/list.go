package handlers

import (
	"fmt"
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
	a.GET("/get", getListHandler)
	a.POST("/update", updateListHandler)
}

func getListHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	list := models.GetList(uid)
	if list == nil {
		models.NewList(uid).Save()
		return utils.Success(res, "")
	}
	gids := []string{}
	for _, gid := range strings.Split(list.GIDs, ";") {
		gids = append(gids, gid)
	}
	return utils.Success(res, models.GetGroupsSlice(gids))
}

func updateListHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	qtype := req.Params["qtype"]
	gid := req.Params["gid"]
	if models.GetGroup(gid) == nil {
		air.ERROR("group not found", utils.M{
			"gid": gid,
			"uid": uid,
		})
		return utils.Error(404, fmt.Errorf("group not found: %s", gid))
	}
	list := models.GetList(uid)
	if list == nil {
		list = models.NewList(uid)
	}
	switch qtype {
	case "add":
		if !strings.Contains(list.GIDs, gid) {
			list.GIDs = strings.Join(
				[]string{list.GIDs, gid}, ";")
		}
		list.Save()
		return utils.Success(res, "")
	case "delete":
		if strings.Contains(list.GIDs, gid) {
			gids := strings.Split(list.GIDs, ";")
			for k, v := range gids {
				if v == gid {
					gids = append(gids[:k], gids[k+1:]...)
					break
				}
			}
			list.GIDs = strings.Join(gids, ";")
		}
		list.Save()
		return utils.Success(res, "")
	}
	return utils.Error(400,
		fmt.Errorf("no specific query type: %s", qtype))
}
