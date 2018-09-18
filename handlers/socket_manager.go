package handlers

import (
	"strings"
	"sync"

	"github.com/BouncyElf/chat/models"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

type SocketManager struct {
	uid       string
	msg       *models.Message
	newMsg    chan struct{}
	shutdown  chan struct{}
	writeLock *sync.Mutex
	sendChan  chan struct{}
}

var (
	mu = &sync.Mutex{}
)

func newSocketManager(uid string) *SocketManager {
	return &SocketManager{
		uid:       uid,
		newMsg:    make(chan struct{}, 1),
		shutdown:  make(chan struct{}),
		writeLock: &sync.Mutex{},
		sendChan:  make(chan struct{}, 1),
	}
}

func (sm *SocketManager) SendMsg(msg *models.Message) {
	sm.sendChan <- struct{}{}
	defer func() {
		<-sm.sendChan
		if sm.msg != nil {
			go sm.msg.Save()
		}
	}()
	if !IsInGroup(sm.uid, msg.To) {
		sm.writeLock.Lock()
		sm.msg = models.NewNotifyMsg(models.UserNotInGroupMsg)
		sm.newMsg <- struct{}{}
		air.ERROR("not in specific group", utils.M{
			"uid":      sm.uid,
			"group id": sm.msg.To,
			"message":  sm.msg,
		})
		return
	}
	group := models.GetGroup(sm.msg.To)
	if group == nil {
		sm.writeLock.Lock()
		sm.msg = models.NewNotifyMsg(models.GroupNotFoundMsg)
		sm.newMsg <- struct{}{}
		air.ERROR("no specific group found", utils.M{
			"uid":      sm.uid,
			"group id": sm.msg.To,
			"message":  sm.msg,
		})
		return
	}
	for _, v := range strings.Split(group.UIDs, ";") {
		if value, ok := users.Get(v); ok {
			me := value.(*SocketManager)
			me.writeLock.Lock()
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
