package handlers

import (
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

func init() {
	a := air.Group{
		Prefix: "/api/user",
	}
	a.POST("/register", registerHandler)
	a.POST("/login", loginHandler)
}

func registerHandler(req *air.Request, res *air.Response) error {
	return utils.Success(res, "")
}

func loginHandler(req *air.Request, res *air.Response) error {
	return utils.Success(res, "")
}
