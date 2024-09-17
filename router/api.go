package router

import (
	"chatapp/handler"
	"chatapp/middleware"
	websocket "chatapp/web_socket"

	"github.com/labstack/echo/v4"
)

type API struct {
	Echo           *echo.Echo
	UserHandler    handler.UserHandler
	ConversHandler handler.ConversHandler
	WsHandler      *websocket.WebSocketHandler
}

func (api *API) SetupRouter() {
	//user
	api.Echo.POST("/user/sign-up", api.UserHandler.HandleSignUp)
	api.Echo.POST("/user/sign-in", api.UserHandler.HandleSignIn)
	friends:= api.Echo.Group("/friends", middleware.JWTMiddleWare())
	friends.POST("/user/friendship", api.UserHandler.HandleFriendship)
	friends.GET("/user/list-friend", api.UserHandler.HandleGetListFriend)
	//chat
	convers := api.Echo.Group("/convers", middleware.JWTMiddleWare())
	convers.POST("/create", api.ConversHandler.CreateConversHandler)
	convers.GET("/list", api.ConversHandler.HandleLoadListConvers)
	convers.GET("/list/mem", api.ConversHandler.HandleLoadListMem)
	convers.POST("/message/send", api.ConversHandler.HandleSendMessage)
	convers.GET("/message/list", api.ConversHandler.HandleLoadMessages)
	convers.PUT("/message/seen", api.ConversHandler.HandleSeenMessage)
		//ws
	api.Echo.GET("/ws/message", api.WsHandler.HandleWebSocketChat)

}
