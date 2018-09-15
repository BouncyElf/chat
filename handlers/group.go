package handlers

import (
	"github.com/BouncyElf/chat/models"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

func init() {
	a := air.Group{
		Prefix: "/api/group",
	}
	a.POST("/new", newGroupHandler)
	a.POST("/list", getGroupListHandler)
	a.POST("/addmem", addMemberHandler)
}

func newGroupHandler(req *air.Request, res *air.Response) error {
	return utils.Success(res, "")
}

func getGroupListHandler(req *air.Request, res *air.Response) error {
	return utils.Success(res, "")
}

func addMemberHandler(req *air.Request, res *air.Response) error {
	return utils.Success(res, "")
}

func getGroupUser(gid string) []models.UserInfo {
	return []models.UserInfo{}
}
