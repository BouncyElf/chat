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
			for _, c := range req.Cookies {
				if c.Name == common.AuthCookieName {
					sid = c.Value
					c.Expires = time.Now().
						Add(7 * 24 * time.Hour)
					cookie = c
					break
				}
			}
			if sid == "" {
				air.ERROR("sid not found in cookie")
				if req.Method == "GET" {
					req.URL.Path = "/login"
					res.StatusCode = 302
					return res.Redirect(req.URL.String())
				}
				return utils.Error(401,
					errors.New("sid not found in cookie"))
			}
			v, ok := User(sid)
			if !ok {
				air.ERROR("sid not found in cache")
				if req.Method == "GET" {
					req.URL.Path = "/login"
					res.StatusCode = 302
					return res.Redirect(req.URL.String())
				}
				return utils.Error(400,
					errors.New("sid not found in cache"))
			}
			req.Params["uid"] = v.Uid
			res.Cookies = append(res.Cookies, cookie)
			return next(req, res)
		}
	}
}

func User(sid string) (models.UserInfo, bool) {
	return models.UserInfo{}, true
}
