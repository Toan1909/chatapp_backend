package repo

import (
	"chatapp/model"
	"chatapp/model/req"
	"context"
)

type ConversRepo interface {
	CreateConvers(c context.Context,req req.ReqCreateConvers) (model.Conversation,error)
	AddMember(c context.Context,memId string,conversId string) (error)
	LoadListConvers(c context.Context,userId string)([]model.Conversation,error)
	LoadListMembers(c context.Context,conversId string)([]model.ConversationMember,error)

	SendMessage(c context.Context, reqSendMsg req.SendMessage) (model.Message, error)
	MarkMessageAsSeen(c context.Context ,messageId,userId string) (model.Message, error)
	LoadMessages(c context.Context, conversId string)([]model.Message,error)
}