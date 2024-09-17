package handler

import (
	my_err "chatapp/err"
	"chatapp/model"
	"chatapp/model/req"
	"chatapp/repo"
	websocket "chatapp/web_socket"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type ConversHandler struct {
	ConversRepo repo.ConversRepo
	WsHandler websocket.WebSocketHandler
}
func (r *ConversHandler) CreateConversHandler(c echo.Context) error{
	req := req.ReqCreateConvers{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	convers, err := r.ConversRepo.CreateConvers(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	for _,mem:= range req.ListMember{
		r.ConversRepo.AddMember(c.Request().Context(),mem,convers.ConversationId)
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Created Conversation thành công",
		Data:       convers,
	})
}
func (r *ConversHandler) HandleSendMessage(c echo.Context) error {
    req := req.SendMessage{}
    if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

    validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

    msg, err := r.ConversRepo.SendMessage(c.Request().Context(), req)
    if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	// Gửi tin nhắn tới WebSocket clients qua kênh Broadcast
    r.WsHandler.Broadcast <- msg
	//
    return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Send Message thành công",
		Data:       msg,
	})
}
func (g *ConversHandler) HandleLoadListConvers(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomclaims)
	listConvers, err := g.ConversRepo.LoadListConvers(c.Request().Context(), claims.UserId)
	if err != nil {
		if err == my_err.ConvsersNotFound {
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Select Convers thanh cong",
		Data:       listConvers,
	})

}
func (g *ConversHandler) HandleLoadListMem(c echo.Context) error {
	req := req.ReqLoadMem{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	listMem, err := g.ConversRepo.LoadListMembers(c.Request().Context(),req.ConversationId)
	if err != nil {
		if err == my_err.MemNotFound{
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Load List Member thanh cong",
		Data:       listMem,
	})

}
func (g *ConversHandler) HandleLoadMessages(c echo.Context) error {
	req := req.ReqGetMessage{}
    if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	listMessage, err := g.ConversRepo.LoadMessages(c.Request().Context(),req.ConversId)
	if err != nil {
		if err == my_err.MessageNotFound {
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Select Bookmarks thành công",
		Data:       listMessage,
	})

}
func (g *ConversHandler) HandleSeenMessage(c echo.Context) error {
	req := req.ReqReadReceipt{}
    if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	message,err := g.ConversRepo.MarkMessageAsSeen(c.Request().Context(),req.MessageId,req.UserId)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	g.WsHandler.Broadcast <- message
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    req.UserId+" seen "+req.MessageId,
		Data:       nil,
	})

}