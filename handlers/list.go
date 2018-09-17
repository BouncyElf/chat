package handlers

import (
	"github.com/BouncyElf/chat/gas"

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
	return nil
}

func updateListHandler(req *air.Request, res *air.Response) error {
	return nil
}
