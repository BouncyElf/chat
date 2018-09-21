package models

import (
	"encoding/json"
	"time"

	"github.com/BouncyElf/chat/common"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

var (
	UserNotInGroupMsg Message
	GroupNotFoundMsg  Message
)

func initMessage() error {
	UserNotInGroupMsg = Message{
		From:    "system",
		MType:   air.WebSocketMessageTypeText,
		Type:    "system",
		Content: common.UserNotInGroupMsg,
	}
	GroupNotFoundMsg = Message{
		From:    "system",
		MType:   air.WebSocketMessageTypeText,
		Type:    "system",
		Content: common.GroupNotFoundMsg,
	}
	return nil
}

type Message struct {
	// uuid
	MID  string `gorm:"column:mid;primary_key" json:"id"`
	From string `gorm:"column:from" json:"from"`

	// gid or uid, when system notify `to` is uid
	To      string `gorm:"column:to" json:"group_id"`
	Type    string `gorm:"column:type" json:"type"`
	Content string `gorm:"column:content" json:"content"`
	Time    string `gorm:"column:time" json:"time"`

	MType air.WebSocketMessageType `gorm:"-" json:"-"`
}

func (Message) TableName() string {
	return "message"
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(*m)
}

func (m *Message) Save() {
	if m.MID == "" {
		m.MID = common.NewSnowFlake()
	}
	err := DB.Save(m).Error
	if err != nil {
		air.ERROR("save message to db error", utils.M{
			"err":     err.Error(),
			"message": m,
		})
	}
}

func NewMsg(from string, t air.WebSocketMessageType, b []byte) *Message {
	m := utils.M{}
	_ = json.Unmarshal(b, &m)
	to, _ := m["group_id"].(string)
	msgType, _ := m["type"].(string)
	content, _ := m["content"].(string)
	return &Message{
		MID:     common.NewSnowFlake(),
		From:    from,
		To:      to,
		Type:    msgType,
		Content: content,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		MType:   t,
	}
}

func NewNotifyMsg(msg Message) *Message {
	msg.MID = common.NewSnowFlake()
	msg.Time = time.Now().Format("15:04:05")
	msg.MType = air.WebSocketMessageTypeText
	return &msg
}
