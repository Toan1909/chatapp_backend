package repo

import (
	"chatapp/model"
	"chatapp/model/req"
	"context"
)

type ConversRepo interface {
	CreateConvers(c context.Context,req req.ReqCreateConvers) (model.Conversation,error)
	AddMember(c context.Context,memId string,conversId string) (error)
	LoadListConversation(c context.Context,userId string)([]model.Conversation,error)
	LoadListMembers(c context.Context,conversId string)([]model.ConversationMember,error)

	SendMessage(c context.Context, reqSendMsg req.SendMessage) (model.Message, error)
	UpdateLastMessageSeen(c context.Context,conversId,messageId,userId string)(model.ResponSeenMessage,error)
	LoadMessages(c context.Context, conversId string)([]model.Message,error)
}