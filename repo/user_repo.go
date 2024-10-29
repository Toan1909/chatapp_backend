package repo

import (
	"chatapp/model"
	"chatapp/model/req"
	"context"
)

type UserRepo interface {
	CheckLogIn(c context.Context, loginReq req.ReqSignIn) (model.User, error)
	SaveUser(c context.Context, user model.User) (model.User, error)
	GetUserInfo(c context.Context, userId string ) (model.User, error)
	LoadListFriend(c context.Context, userId string) ([]model.User,error)
	LoadListPending(c context.Context, userId string) ([]model.User,error)
	SearchUser(c context.Context, email,phone string) (model.User,error)
	CheckFriend(c context.Context, userId,friendId string) (string,error)
	SaveFriendShip(c context.Context, userId ,friendId string)(model.FriendShip,error)
	AcceptFriendShip(c context.Context, userId ,friendId string)(error)
}