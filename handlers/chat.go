package handlers

import (
	"strings"

	"github.com/BouncyElf/chat/gas"
	"github.com/BouncyElf/chat/models"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
	cmap "github.com/orcaman/concurrent-map"
)

var (
	users = cmap.New()
)

func init() {
	air.GET("/socket", socketHandler, gas.Auth)
}

func SendMsg(sm *SocketManager, msg *models.Message) {
	if msg == nil {
		air.ERROR("nil message", utils.M{
			"sm.uid": sm.uid,
			"sm.msg": sm.msg,
		})
		return
	}
	var (
		updateUnread    = false
		updateUnreadUID = ""
		updateUnreadGID = ""
	)
	defer func() {
		realMsg := msg
		if sm.msg != nil {
			realMsg = sm.msg
		}
		realMsg.Save()
		if updateUnread &&
			!hasUnreadMsg(updateUnreadUID, updateUnreadGID) {
			updateUnreadMsg(
				updateUnreadUID,
				updateUnreadGID,
				realMsg.MID)
		}
	}()
	if sm == nil {
		// system notify
		if msg.Time == "" {
			msg = models.NewNotifyMsg(*msg)
		}
		to, ok := users.Get(msg.To)
		if ok {
			me := to.(*SocketManager)
			me.writeChan <- struct{}{}
			me.msg = msg
			me.newMsg <- struct{}{}
		} else {
			updateUnread = true
			updateUnreadUID = msg.To
			updateUnreadGID = msg.To
		}
		return
	}
	if !IsInGroup(sm.uid, msg.To) {
		sm.writeChan <- struct{}{}
		sm.msg = models.NewNotifyMsg(models.UserNotInGroupMsg)
		sm.newMsg <- struct{}{}
		air.ERROR("not in specific group", utils.M{
			"uid":      sm.uid,
			"group id": sm.msg.To,
			"message":  sm.msg,
		})
		return
	}
	group := models.GetGroup(msg.To)
	// IsInGroup has already judge if group is nil
	// so, here group can't be nil
	for _, v := range strings.Split(group.UIDs, ";") {
		if value, ok := users.Get(v); ok {
			me := value.(*SocketManager)
			me.writeChan <- struct{}{}
			me.msg = msg
			me.newMsg <- struct{}{}
		} else {
			updateUnread = true
			updateUnreadUID = v
			updateUnreadGID = msg.To
		}
	}
}

func socketHandler(req *air.Request, res *air.Response) error {
	c, err := res.UpgradeToWebSocket()
	if err != nil {
		air.ERROR("upgrade to websocket error", utils.M{
			"request": req,
			"res":     res,
			"err":     err.Error(),
		})
		return utils.Error(500, err)
	}
	defer c.Close()

	me := newSocketManager(req.Params["uid"])
	myInfo := models.GetUserInfo(req.Params["uid"])
	if myInfo == nil {
		air.ERROR("get user info error", utils.M{
			"request": req,
		})
		return utils.Error(500, err)
	}
	users.Set(me.uid, me)

	waitChan := make(chan struct{}, 1)

	go func() {
		for {
			if t, b, err := c.ReadMessage(); err == nil {
				switch t {
				case air.WebSocketMessageTypeText:
					waitChan <- struct{}{}
					go func() {
						SendMsg(
							me,
							models.NewMsg(
								me.uid,
								myInfo.Name,
								t,
								b,
							),
						)
						<-waitChan
					}()
				case air.WebSocketMessageTypeConnectionClose:
					me.Close()
					return
				}
			} else {
				air.ERROR("socket msg error", utils.M{
					"type":    t,
					"err":     err.Error(),
					"content": string(b),
				})
				me.Close()
				return
			}
		}
	}()

	for {
		select {
		case <-me.newMsg:
			if v, err := me.msg.Marshal(); err == nil {
				err = c.WriteMessage(me.msg.MType, v)
				if err != nil {
					air.ERROR("send socket msg error",
						utils.M{
							"content": me.msg,
							"to":      me.uid,
							"err":     err,
						})
					me.Close()
				}
				<-me.writeChan
			}
		case <-me.shutdown:
			break
		}
	}

	return nil
}
