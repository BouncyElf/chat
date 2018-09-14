package handlers

import (
	"errors"
	"sync"

	"github.com/BouncyElf/chat/gas"
	"github.com/BouncyElf/chat/utils"
	"github.com/aofei/air"
)

var (
	users map[string]*SocketManager
	mu    = &sync.Mutex{}
)

func init() {
	air.GET("/socket", socketHandler, gas.Auth)
}

func socketHandler(req *air.Request, res *air.Response) error {
	c, err := res.UpgradeToWebSocket()
	if err != nil {
		air.ERROR("upgrade to websocket error", utils.M{
			"request": req,
			"error":   err.Error(),
		})
		return utils.Error(500, err)
	}
	defer c.Close()

	name := req.Params["name"]
	if _, ok := users[name]; ok {
		air.ERROR("duplicate name", utils.M{
			"req":  req,
			"name": name,
		})
		return utils.Error(400, errors.New("duplicate name"))
	}

	me := newSocketManager(name)
	if users == nil {
		users = make(map[string]*SocketManager)
	}
	users[name] = me

	go func() {
		for {
			if t, b, err := c.ReadMessage(); err == nil {
				switch t {
				case air.WebSocketMessageTypeText:
					mu.Lock()
					me.SendMsg(newMsg(name, t, b))
					mu.Unlock()
				case air.WebSocketMessageTypeBinary:
				case air.WebSocketMessageTypeConnectionClose:
					delete(users, name)
					me.Close()
					return
				}
			} else {
				air.ERROR("socket msg error", utils.M{
					"type":    t,
					"err":     err.Error(),
					"content": string(b),
				})
				delete(users, name)
				me.Close()
				return
			}
		}
	}()

	for {
		select {
		case <-me.newMsg:
			if v, err := me.m.Marshal(); err == nil {
				err = c.WriteMessage(me.m.Mtype, v)
				if err != nil {
					air.ERROR("send socket msg error",
						utils.M{
							"content": me.m,
							"to":      me.name,
							"err":     err,
						})
				}
				me.rwLock.Unlock()
			}
		case <-me.shutdown:
			break
		}
	}

	return nil
}
