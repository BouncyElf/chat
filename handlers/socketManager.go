package handlers

import (
	"sync"
)

type SocketManager struct {
	name     string
	m        *Message
	newMsg   chan struct{}
	shutdown chan struct{}
	rwLock   *sync.Mutex
}

func newSocketManager(name string) *SocketManager {
	return &SocketManager{
		name:     name,
		newMsg:   make(chan struct{}, 1),
		shutdown: make(chan struct{}),
		rwLock:   &sync.Mutex{},
	}
}

func (sm *SocketManager) SendMsg(m *Message) {
	for _, v := range users {
		v.rwLock.Lock()
		v.m = m
		v.newMsg <- struct{}{}
	}
}

func (sm *SocketManager) Close() {
	sm.shutdown <- struct{}{}
}

func CloseSocket() {
	for _, v := range users {
		v := v
		go func() {
			v.Close()
		}()
	}
}
