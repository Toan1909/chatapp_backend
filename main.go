package main

import (
	"chatapp/db"
	"chatapp/handler"
	repoimpl "chatapp/repo/repo-impl"
	"chatapp/router"
	websocket "chatapp/web_socket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
    sql := &db.Sql{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "190901",
		Dbname:   "chat_app",
	}
    sql.ConnectDb()
    defer sql.CloseDb()
    e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	wsHandler := websocket.NewWebSocketHandler()
    userHandler := handler.UserHandler{
        UserRepo: repoimpl.NewUserRepo(sql),
    }
	conversHandler := handler.ConversHandler{
		ConversRepo: repoimpl.NewConversRepoImpl(sql),
		WsHandler: *wsHandler,
	}
	
    api :=router.API{
        Echo: e,
        UserHandler: userHandler,
		ConversHandler: conversHandler,
		WsHandler: wsHandler,
    }

	// Khởi chạy goroutine để lắng nghe và broadcast tin nhắn
	go wsHandler.BroadcastMessages()
	//================================
    api.SetupRouter()
    e.Logger.Fatal(e.Start(":3000"))
}