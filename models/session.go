package models

import (
	"github.com/BouncyElf/chat/common"
)

type Session struct {
	// uuid
	SID  string
	Info UserInfo
}

var GlobalSession = map[string]*Session{}

func (s *Session) Save() {
	GlobalSession[s.SID] = s
}

func NewSession(info *UserInfo) *Session {
	return &Session{
		SID:  common.NewSnowFlake(),
		Info: *info,
	}
}

func GetSession(sid string) *Session {
	return GlobalSession[sid]
}
