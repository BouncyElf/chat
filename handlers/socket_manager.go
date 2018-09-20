package handlers

import (
	"sync"

	"github.com/BouncyElf/chat/models"
)

type SocketManager struct {
	uid       string
	msg       *models.Message
	newMsg    chan struct{}
	shutdown  chan struct{}
	writeChan chan struct{}
}

var (
	mu = &sync.Mutex{}
)

func newSocketManager(uid string) *SocketManager {
	return &SocketManager{
		uid:       uid,
		newMsg:    make(chan struct{}, 1),
		shutdown:  make(chan struct{}),
		writeChan: make(chan struct{}, 1),
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
