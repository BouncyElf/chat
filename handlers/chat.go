package handlers

import (
	"sync"

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

	me := newSocketManager(req.Params["uid"])
	users.Set(me.uid, me)
	mu := &sync.Mutex{}

	go func() {
		for {
			if t, b, err := c.ReadMessage(); err == nil {
				switch t {
				case air.WebSocketMessageTypeText:
					mu.Lock()
					me.SendMsg(models.NewMsg(me.uid, t, b))
					mu.Unlock()
				case air.WebSocketMessageTypeBinary:
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
				me.mu.Unlock()
			}
		case <-me.shutdown:
			break
		}
	}

	return nil
}
