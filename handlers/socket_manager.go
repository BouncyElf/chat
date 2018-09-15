package handlers

import (
	"sync"
)

type SocketManager struct {
	uid      string
	msg      *Message
	newMsg   chan struct{}
	shutdown chan struct{}
	mu       *sync.Mutex
}

func newSocketManager(uid string) *SocketManager {
	return &SocketManager{
		uid:      uid,
		newMsg:   make(chan struct{}, 1),
		shutdown: make(chan struct{}),
		mu:       &sync.Mutex{},
	}
}

func (sm *SocketManager) SendMsg(msg *Message) {
	for _, v := range getGroupUser(sm.msg.To) {
		// TODO: save to db
		if value, ok := users.Get(v.Uid); ok {
			me := value.(*SocketManager)
			me.mu.Lock()
			me.msg = msg
			me.newMsg <- struct{}{}
		}
	}
}

func (sm *SocketManager) Close() {
	sm.shutdown <- struct{}{}
	users.Remove(sm.uid)
}

func CloseSocket() {
	for _, v := range users.Keys() {
		if value, ok := users.Get(v); ok {
			sm, _ := value.(*SocketManager)
			sm.Close()
		}
	}
}
