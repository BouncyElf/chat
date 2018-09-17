package gas

import (
	"errors"
	"time"

	"github.com/BouncyElf/chat/common"
	"github.com/BouncyElf/chat/models"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

func AuthHandler() air.Gas {
	return func(next air.Handler) air.Handler {
		return func(req *air.Request, res *air.Response) error {
			// if air.DebugMode {
			// 	air.INFO("debug mode, pass")
			// 	return next(req, res)
			// }
			sid := ""
			cookie := &air.Cookie{}
			c, ok := req.Cookies[common.AuthCookieName]
			if !ok {
				air.ERROR("sid not found in cookie")
				return utils.Error(401,
					errors.New("sid not found"))
			}
			sid = c.Value
			c.Expires = time.Now().Add(7 * 24 * time.Hour)
			cookie = c
			v := models.GetSession(sid)
			if v == nil {
				air.ERROR("sid not found in session")
				return utils.Error(401,
					errors.New("session expired"))
			}
			req.Params["uid"] = v.Info.UID
			res.Cookies[common.AuthCookieName] = cookie
			return next(req, res)
		}
	}
}
