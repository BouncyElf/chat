package handlers

import (
	"encoding/json"
	"time"

	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

type Message struct {
	From    string                   `json:"from"`
	To      string                   `json:"group_id"`
	MType   air.WebSocketMessageType `json:"-"`
	Type    string                   `json:"type"`
	Content string                   `json:"content"`
	Time    string                   `json:"time"`
}

func newMsg(from string, t air.WebSocketMessageType, b []byte) *Message {
	m := utils.M{}
	_ = json.Unmarshal(b, &m)
	gid, _ := m["group_id"].(string)
	msgType, _ := m["type"].(string)
	content, _ := m["content"].(string)
	return &Message{
		From:    from,
		To:      gid,
		Type:    msgType,
		MType:   t,
		Content: content,
		Time:    time.Now().Format("15:04:05"),
	}
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(*m)
}
