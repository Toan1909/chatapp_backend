package websocket

import (
	"chatapp/model"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	Clients   map[string]map[*websocket.Conn]bool // Map conversation_id -> (map client_conn -> boolean)
	Broadcast chan model.Message
	Upgrader  websocket.Upgrader // WebSocket upgrader
}

// NewWebSocketHandler khởi tạo WebSocketHandler mới
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		Clients:   make(map[string]map[*websocket.Conn]bool),
		Broadcast: make(chan model.Message),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleWebSocketChat xử lý kết nối WebSocket
func (h *WebSocketHandler) HandleWebSocketChat(c echo.Context) error {
	conversationID := c.QueryParam("conversationId")
	if conversationID == "" {
		return c.JSON(http.StatusBadRequest, "Missing conversation_id")
	}

	ws, err := h.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return err
	}
	defer ws.Close()

	if _, ok := h.Clients[conversationID]; !ok {
		h.Clients[conversationID] = make(map[*websocket.Conn]bool)
	}
	h.Clients[conversationID][ws] = true

	for {
		var msg model.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading json:", err)
			delete(h.Clients[conversationID], ws)
			if len(h.Clients[conversationID]) == 0 {
				delete(h.Clients, conversationID)
			}
			break
		}

		h.Broadcast <- msg
	}

	return nil
}

// BroadcastMessages lắng nghe kênh broadcast và gửi tin nhắn tới tất cả các client
func (h *WebSocketHandler) BroadcastMessages() {
	for {
		msg := <-h.Broadcast
		if clients, ok := h.Clients[msg.ConversationId]; ok {
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Println("Error writing json:", err)
					client.Close()
					delete(clients, client)
					if len(clients) == 0 {
						delete(h.Clients, msg.ConversationId)
					}
				}
			}
		}
	}
}
