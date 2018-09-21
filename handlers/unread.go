package handlers

import (
	"github.com/BouncyElf/chat/gas"
	"github.com/BouncyElf/chat/models"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

func init() {
	a := air.Group{
		Prefix: "/api/unread",
		Gases:  []air.Gas{gas.Auth},
	}
	a.POST("/update", updateUnreadHandler)
}

func updateUnreadHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	gid := req.Params["gid"]
	lastMID := req.Params["last_mid"]
	unread := models.GetUnread(uid, gid)
	if unread == nil {
		unread = models.NewUnread(uid, gid, lastMID)
	}
	unread.LastMID = lastMID
	unread.Count++
	go unread.Save()
	return utils.Success(res, "")
}
