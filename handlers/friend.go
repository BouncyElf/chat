package handlers

import (
	"errors"
	"strings"

	"github.com/BouncyElf/chat/common"
	"github.com/BouncyElf/chat/gas"
	"github.com/BouncyElf/chat/models"
	"github.com/BouncyElf/chat/utils"

	"github.com/aofei/air"
)

func init() {
	a := air.Group{
		Prefix: "/api/friend",
		Gases:  []air.Gas{gas.Auth},
	}
	a.POST("/list", listFriendHandler)
	a.POST("/confirm/friend", confirmAddFriend)
	a.POST("/add", addFriendHandler)
	a.POST("/delete", deleteFriendHandler)
}

func listFriendHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	friends := models.GetFriends(uid)
	friendsInfos := models.GetUserInfosSlice(friends)
	lists := []interface{}{}
	for _, friend := range friendsInfos {
		group := models.GetPrivateChatGroup(uid, friend.UID)
		if group != nil {
			lists = append(lists, utils.M{
				"info": friend,
				"gid":  group.GID,
			})
		}
	}
	return utils.Success(res, lists)
}

func addFriendHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	displayID := req.Params["display_id"]
	info := models.GetUserInfoByDisplayID(displayID)
	if info == nil {
		air.ERROR("userinfo not found", utils.M{
			"display_id": displayID,
		})
		return utils.Error(404, errors.New("userinfo not found"))
	}
	myInfo := models.GetUserInfo(uid)
	if myInfo == nil {
		air.ERROR("userinfo not found", utils.M{
			"uid": uid,
		})
		return utils.Error(500, errors.New("server internal error"))
	}
	SendMsg(nil, &models.Message{
		From:    myInfo.UID,
		To:      info.UID,
		Type:    common.MsgTypeSystem + "/" + "friend",
		Content: myInfo.Name + "申请成为您的好友",
	})
	return utils.Success(res, "")
}

func confirmAddFriend(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	tuid := req.Params["tuid"]
	relation := models.GetRelations([]string{uid, tuid})
	if relation[uid] == nil {
		relation[uid] = models.NewRelation(uid)
	}
	if !strings.Contains(relation[uid].UIDs, tuid) {
		relation[uid].UIDs = strings.Join([]string{
			relation[uid].UIDs,
			tuid,
		}, ";")
	}
	go relation[uid].Save()
	if relation[tuid] == nil {
		relation[tuid] = models.NewRelation(tuid)
	}
	if !strings.Contains(relation[tuid].UIDs, uid) {
		relation[tuid].UIDs = strings.Join([]string{
			relation[tuid].UIDs,
			uid,
		}, ";")
	}
	go relation[tuid].Save()
	info := models.GetUserInfos([]string{uid, tuid})
	if info[uid] != nil && info[tuid] != nil {
		name := strings.Join(
			[]string{
				info[uid].Name,
				info[tuid].Name,
			},
			";",
		)
		go models.NewGroup(
			[]string{
				uid,
				tuid,
			},
			common.ChatTypePrivate,
			name,
		).Save()
	}
	return utils.Success(res, "")
}

func deleteFriendHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	tuid := req.Params["tuid"]
	relations := models.GetRelations([]string{uid, tuid})
	if relations[uid] != nil {
		if strings.Contains(relations[uid].UIDs, tuid) {
			friends := strings.Split(relations[uid].UIDs, ";")
			for k, v := range friends {
				if v == tuid {
					friends = append(friends[:k],
						friends[k+1:]...)
					break
				}
			}
			relations[uid].UIDs = strings.Join(friends, ";")
			go relations[uid].Save()
		}
	}
	if relations[tuid] != nil {
		if strings.Contains(relations[tuid].UIDs, uid) {
			friends := strings.Split(relations[tuid].UIDs, ";")
			for k, v := range friends {
				if v == uid {
					friends = append(friends[:k],
						friends[k+1:]...)
					break
				}
			}
			relations[tuid].UIDs = strings.Join(friends, ";")
			go relations[tuid].Save()
		}
	}
	return utils.Success(res, "")
}
