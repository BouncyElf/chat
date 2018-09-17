package handlers

import (
	"errors"
	"fmt"

	"github.com/BouncyElf/chat/common"
	"github.com/BouncyElf/chat/models"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

func init() {
	a := air.Group{
		Prefix: "/api/user",
	}

	a.POST("/register", registerHandler)
	a.POST("/login", loginHandler)
}

func registerHandler(req *air.Request, res *air.Response) error {
	name := req.Params["name"]
	username := req.Params["username"]
	password := req.Params["password"]
	if name == "" || username == "" || password == "" {
		return utils.Error(400, errors.New("param empty"))
	}

	u := models.GetUserByUsername(username)
	if u != nil {
		air.ERROR("register with duplicate username", utils.M{
			"username": username,
			"name":     name,
		})
		return utils.Error(409,
			fmt.Errorf("duplicate username: %s", username))
	}

	user := models.NewUser(username, utils.MD5(password))
	sid, err := newUser(user, name)
	if err != nil {
		air.ERROR("new user error", utils.M{
			"err":  err.Error(),
			"user": user,
		})
		return utils.Error(500, err)
	}
	res.Cookies[common.AuthCookieName] = common.NewAuthCookie(sid)
	return utils.Success(res, "")
}

func newUser(u *models.User, name string) (string, error) {
	if u == nil {
		return "", errors.New("user is nil")
	}
	if u.UID == "" {
		u.UID = common.NewUUID()
	}
	userInfo := models.NewUserInfo(u.UID, name)
	list := models.NewList(u.UID)
	relation := models.NewRelation(u.UID)
	session := models.NewSession(userInfo)
	u.Save()
	userInfo.Save()
	list.Save()
	relation.Save()
	session.Save()
	return session.SID, nil
}

func loginHandler(req *air.Request, res *air.Response) error {
	username := req.Params["username"]
	password := req.Params["password"]
	u := models.GetUserByUsername(username)
	if u == nil {
		air.ERROR("user not found", utils.M{
			"username": username,
		})
		return utils.Error(404, errors.New("user not found"))
	}
	pwd := utils.MD5(password)
	if pwd != password {
		air.ERROR("wrong password", utils.M{
			"username": username,
			"host":     req.ClientIP.String(),
		})
		return utils.Error(400, errors.New("wrong password"))
	}
	userInfo := models.GetUserInfo(u.UID)
	if userInfo == nil {
		userInfo = models.NewUserInfo(u.UID, "name")
		userInfo.Save()
	}
	session := models.NewSession(userInfo)
	session.Save()
	res.Cookies[common.AuthCookieName] = common.NewAuthCookie(session.SID)
	return utils.Success(res, "")
}
