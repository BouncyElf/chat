package handlers

import (
	"github.com/BouncyElf/chat/gas"
	"github.com/aofei/air"
)

func init() {
	a := air.Group{
		Prefix: "/api/friend",
		Gases:  []air.Gas{gas.Auth},
	}
	a.GET("/list", listFriendHandler)
	a.POST("/add", addFriendHandler)
	a.POST("/delete", deleteFriendHandler)
}

func listFriendHandler(req *air.Request, res *air.Response) error {
	return nil
}

func addFriendHandler(req *air.Request, res *air.Response) error {
	return nil
}

func deleteFriendHandler(req *air.Request, res *air.Response) error {
	return nil
}
