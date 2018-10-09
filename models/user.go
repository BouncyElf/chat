package models

import (
	"github.com/BouncyElf/chat/common"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

type User struct {
	UID      string `gorm:"column:uid;primary_key"`
	Username string `gorm:"column:username;unique"`
	Password string `gorm:"column:password;not null"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) Save() {
	if u.UID == "" {
		u.UID = common.NewSnowFlake()
	}
	err := DB.Save(u).Error
	if err != nil {
		air.ERROR("save user to db error", utils.M{
			"err":  err.Error(),
			"user": u,
		})
	}
}

type UserInfo struct {
	// uuid
	UID string `gorm:"column:uid;primary_key" json:"uid"`

	// uid + name uuid
	DisplayID string `gorm:"column:display_id;unique" json:"display_id"`
	Name      string `gorm:"column:name;unique" json:"name"`
	Bio       string `gorm:"column:bio" json:"bio"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

func (info *UserInfo) Save() {
	err := DB.Save(info).Error
	if err != nil {
		air.ERROR("save userinfo to db error", utils.M{
			"err":       err.Error(),
			"user info": info,
		})
	}
}

func NewUser(uname string, pwd string) *User {
	return &User{
		UID:      common.NewSnowFlake(),
		Username: uname,
		Password: pwd,
	}
}

func GetUser(uid string) *User {
	res := &User{}
	err := DB.Where("uid = ?", uid).Find(res).Error
	if err != nil {
		air.ERROR("get user from db error", utils.M{
			"err": err.Error(),
			"uid": uid,
		})
		return nil
	}
	return res
}

func GetUserByUsername(username string) *User {
	res := &User{}
	err := DB.Where("username = ?", username).Find(res).Error
	if err != nil {
		air.ERROR("get user from db error", utils.M{
			"err":      err.Error(),
			"username": username,
		})
		return nil
	}
	return res
}

func NewUserInfo(uid string, name string) *UserInfo {
	return &UserInfo{
		UID:       uid,
		DisplayID: common.NewSnowFlake(),
		Name:      name,
		Bio:       "",
	}
}

func GetUserInfo(uid string) *UserInfo {
	res := &UserInfo{}
	err := DB.Where("uid = ?", uid).Find(res).Error
	if err != nil {
		air.ERROR("get userinfo from db error", utils.M{
			"err": err.Error(),
			"uid": uid,
		})
		return nil
	}
	return res
}

func GetUserInfoByDisplayID(displayID string) *UserInfo {
	res := &UserInfo{}
	err := DB.Where("display_id = ?", displayID).Find(res).Error
	if err != nil {
		air.ERROR("get userinfo from db error", utils.M{
			"err":        err.Error(),
			"display_id": displayID,
		})
		return nil
	}
	return res
}

func GetUserInfos(uids []string) map[string]*UserInfo {
	userInfos := GetUserInfosSlice(uids)
	res := map[string]*UserInfo{}
	for _, v := range userInfos {
		res[v.UID] = v
	}
	return res
}

func GetUserInfosSlice(uids []string) []*UserInfo {
	userInfos := []*UserInfo{}
	err := DB.Where("uid in (?)", uids).Find(&userInfos).Error
	if err != nil {
		air.ERROR("get user infos error", utils.M{
			"err":  err.Error(),
			"uids": uids,
		})
	}
	return userInfos
}
