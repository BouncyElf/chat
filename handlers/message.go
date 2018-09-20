package handlers

import (
	"github.com/BouncyElf/chat/gas"

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
	return nil
}

// recentMsgHandler return specific group's recent message.
func recentMsgHandler(req *air.Request, res *air.Response) error {
	return nil
}
