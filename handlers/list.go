package handlers

import (
	"fmt"
	"strconv"
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
	gids := []int64{}
	for _, v := range strings.Split(list.GIDs, ";") {
		gid, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			air.ERROR("parse gid string to int64 error", utils.M{
				"err": err.Error(),
				"gid": v,
			})
			continue
		}
		gids = append(gids, gid)
	}
	return utils.Success(res, models.GetGroupsSlice(gids))
}

func updateListHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	qtype := req.Params["qtype"]
	gidstring := req.Params["gid"]
	gid, err := utils.ParseInt64(req, "gid")
	if err != nil {
		air.ERROR("invalid gid", utils.M{
			"gid": gidstring,
			"uid": uid,
		})
		return utils.Error(404,
			fmt.Errorf("invalid gid: %s", gidstring))
	}
	if models.GetGroup(gid) == nil {
		air.ERROR("group not found", utils.M{
			"gid": gid,
			"uid": uid,
		})
		return utils.Error(404, fmt.Errorf("group not found: %d", gid))
	}
	list := models.GetList(uid)
	if list == nil {
		list = models.NewList(uid)
	}
	gids := strings.Split(list.GIDs, ";")
	switch qtype {
	case "add":
		for _, v := range gids {
			if v == gidstring {
				list.Save()
				return utils.Success(res, "")
			}
		}
		list.GIDs = strings.Join(append(gids, gidstring), ";")
		list.Save()
		return utils.Success(res, "")
	case "delete":
		for k, v := range gids {
			if v == gidstring {
				gids = append(gids[:k], gids[k+1:]...)
				break
			}
		}
		list.GIDs = strings.Join(gids, ";")
		list.Save()
		return utils.Success(res, "")
	}
	return utils.Error(400,
		fmt.Errorf("no specific query type: %s", qtype))
}
