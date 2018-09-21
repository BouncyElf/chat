package models

import "github.com/BouncyElf/chat/common"

type Session struct {
	// uuid
	SID  string
	Info UserInfo
}

func (s *Session) Save() {
}

func NewSession(info *UserInfo) *Session {
	return &Session{
		SID:  common.NewSnowFlake(),
		Info: *info,
	}
}

func GetSession(sid string) *Session {
	return nil
}
