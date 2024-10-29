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
	api.Echo.GET("/user/profile", api.UserHandler.HandleGetProfile)
	friends := api.Echo.Group("/friends", middleware.JWTMiddleWare())
	friends.POST("/friendship", api.UserHandler.HandleFriendship)
	friends.PUT("/accept", api.UserHandler.HandleAcceptFriendship)
	friends.POST("/list-friend", api.UserHandler.HandleGetListFriend)
	friends.POST("/list-pending", api.UserHandler.HandleGetListPending)
	friends.POST("/check-friend", api.UserHandler.HandleCheckFriend)
	friends.GET("/profile", api.UserHandler.HandleGetProfile)
	friends.POST("/search", api.UserHandler.HandleSearchUser)
	//chat
	convers := api.Echo.Group("/convers", middleware.JWTMiddleWare())
	convers.POST("/create", api.ConversHandler.CreateConversHandler)
	convers.GET("/list", api.ConversHandler.HandleLoadListConvers)
	convers.POST("/list/mem", api.ConversHandler.HandleLoadListMem)
	convers.POST("/message/send", api.ConversHandler.HandleSendMessage)
	convers.POST("/message/seen", api.ConversHandler.HandleSeenMessage)
	convers.POST("/message/list", api.ConversHandler.HandleLoadMessages)
	//convers.PUT("/message/seen", api.ConversHandler.HandleSeenMessage)
	//ws
	api.Echo.GET("/ws/message", api.WsHandler.HandleWebSocket)

}
