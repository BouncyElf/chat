package handlers

import (
	"strings"

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
	a.GET("/list", listFriendHandler)
	a.POST("/add", addFriendHandler)
	a.POST("/delete", deleteFriendHandler)
}

func listFriendHandler(req *air.Request, res *air.Response) error {
	friends := models.GetFriends(req.Params["uid"])
	return utils.Success(res, models.GetUserInfosSlice(friends))
}

func addFriendHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	tuid := req.Params["tuid"]
	_, _ = uid, tuid
	// TODO: send system msg to confirm or abort
	return utils.Success(res, "")
}

func confirmAddFriend(uid, tuid string) error {
	// TODO: update uid and tuid's relation
	// add group of uid and tuid, ChatTypePrivate
	return nil
}

func deleteFriendHandler(req *air.Request, res *air.Response) error {
	uid := req.Params["uid"]
	tuid := req.Params["tuid"]
	relations := models.GetRelations([]string{uid, tuid})
	if relations[uid] != nil {
		friends := strings.Split(relations[uid].UIDs, ";")
		for k, v := range friends {
			if v == tuid {
				friends = append(friends[:k], friends[k+1:]...)
				break
			}
		}
		relations[uid].UIDs = strings.Join(friends, ";")
		relations[uid].Save()
	}
	if relations[tuid] != nil {
		friends := strings.Split(relations[tuid].UIDs, ";")
		for k, v := range friends {
			if v == uid {
				friends = append(friends[:k], friends[k+1:]...)
				break
			}
		}
		relations[tuid].UIDs = strings.Join(friends, ";")
		relations[tuid].Save()
	}
	// TODO: delete uid, tuid list
	return utils.Success(res, "")
}
