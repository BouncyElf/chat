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
		Prefix: "/api/message",
		Gases:  []air.Gas{gas.Auth},
	}
	a.POST("/list", listMsgHandler)
	a.POST("/recent", recentMsgHandler)
}

// listMsgHandler return the list page recent message.
func listMsgHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	gids := strings.Split(req.Params["gids"], ";")
	data := utils.M{}
	for _, gid := range gids {
		lastMsg := models.GetLastMsg(gid)
		if lastMsg == nil {
			break
		}
		data[gid] = utils.M{
			"show_msg": lastMsg,
			"unread":   hasUnreadMsg(uid, gid),
		}
	}
	return utils.Success(res, data)
}

// recentMsgHandler return specific group's recent message.
func recentMsgHandler(req *air.Request, res *air.Response) error {
	gid := req.Params["gid"]
	count, _ := utils.ParseInt(req, "count")
	minMID := req.Params["min_mid"]
	if count == 0 {
		count = 10
	}
	return utils.Success(res, models.GetMessagesSlice(gid, minMID, count))
}
