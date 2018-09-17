package common

import (
	"github.com/aofei/air"

	"github.com/bwmarrin/snowflake"
)

var SnowFlakeNode *snowflake.Node

func InitCommon() {
	SnowFlakeNode, _ = snowflake.NewNode(0)
}

func NewAuthCookie(sid string) *air.Cookie {
	return &air.Cookie{
		Name:  AuthCookieName,
		Value: sid,
		Path:  "/",
	}
}
