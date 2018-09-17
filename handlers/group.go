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
		Prefix: "/api/group",
		Gases:  []air.Gas{gas.Auth},
	}
	a.POST("/new", newGroupHandler)
	a.POST("/addmem", addMemberHandler)
}

func newGroupHandler(req *air.Request, res *air.Response) error {
	return utils.Success(res, "")
}

func addMemberHandler(req *air.Request, res *air.Response) error {
	return utils.Success(res, "")
}

func IsInGroup(uid string, gid int64) bool {
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
