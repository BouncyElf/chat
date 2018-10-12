package handlers

import (
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
	a.POST("/recent", recentMsgHandler)
	a.POST("/system", systemMsgHandler)
}

// recentMsgHandler return specific group's recent message.
func recentMsgHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	gid := req.Params["gid"]
	count, _ := utils.ParseInt(req, "count")
	minMID := req.Params["min_mid"]
	if count == 0 {
		count = 10
	}
	if minMID == "" {
		if hasUnreadMsg(uid, gid) {
			unread := models.GetUnread(uid, gid)
			if unread != nil {
				minMID = unread.LastMID
			} else {
				minMID = "0"
			}
		} else {
			minMID = "0"
		}
	}
	return utils.Success(res, models.GetMessagesSlice(gid, minMID, count))
}

func systemMsgHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	gid := uid
	minMID := ""
	count := 99
	isUnread := hasUnreadMsg(uid, gid)
	if isUnread {
		unread := models.GetUnread(uid, gid)
		if unread != nil {
			minMID = unread.LastMID
		} else {
			minMID = "0"
		}
	} else {
		minMID = "0"
	}
	notifyMsg := models.GetMessagesSlice(gid, minMID, count)
	air.INFO("notify msg", utils.M{
		"msg":    notifyMsg,
		"gid":    gid,
		"minMID": minMID,
		"count":  count,
	})
	return utils.Success(res, utils.M{
		"unread": isUnread,
		"msg":    notifyMsg,
	})
}
