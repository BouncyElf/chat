package gas

import (
	"errors"
	"time"

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
			name := ""
			cookie := &air.Cookie{}
			for _, c := range req.Cookies {
				if c.Name == "name" {
					name = c.Value
					c.Expires = time.Now().
						Add(7 * 24 * time.Hour)
					cookie = c
					break
				}
			}
			if name == "" {
				air.ERROR("name not found in cookie")
				if req.Method == "GET" {
					req.URL.Path = "/"
					res.StatusCode = 302
					return res.Redirect(req.URL.String())
				}
				return utils.Error(401,
					errors.New("name not found in cookie"))
			}
			v, ok := models.GetUser(name)
			if !ok {
				air.ERROR("name not found in cache")
				if req.Method == "GET" {
					req.URL.Path = "/"
					res.StatusCode = 302
					return res.Redirect(req.URL.String())
				}
				return utils.Error(400,
					errors.New("name not found in cache"))
			}
			req.Params["name"] = v.Name
			res.Cookies = append(res.Cookies, cookie)
			return next(req, res)
		}
	}
}
