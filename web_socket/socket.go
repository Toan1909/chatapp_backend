package websocket

import (
	"chatapp/model"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	Clients   map[string]map[*websocket.Conn]bool // Map user_id -> (map client_conn -> boolean)
	Broadcast chan model.ResponseWs
	Upgrader  websocket.Upgrader // WebSocket upgrader
}

// NewWebSocketHandler khởi tạo WebSocketHandler mới
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		Clients:   make(map[string]map[*websocket.Conn]bool),
		Broadcast: make(chan model.ResponseWs),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleWebSocket xử lý kết nối WebSocket
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	wsId := c.QueryParam("wsId")
	if wsId == "" {
		return c.JSON(http.StatusBadRequest, "Missing wsId")
	}

	ws, err := h.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return err
	}
	defer ws.Close()

	if _, ok := h.Clients[wsId]; !ok {
		h.Clients[wsId] = make(map[*websocket.Conn]bool)
	}
	h.Clients[wsId][ws] = true

	for {
		_, _, err := ws.NextReader() // Chỉ lắng nghe kết nối mà không đọc dữ liệu
		if err != nil {
			log.Println("Client disconnected or error:", err)

			// Khi có lỗi, xóa kết nối khỏi danh sách
			delete(h.Clients[wsId], ws)

			// Nếu không còn client nào trong wsId, xóa luôn wsId
			if len(h.Clients[wsId]) == 0 {
				delete(h.Clients, wsId)
			}

			// Thoát khỏi vòng lặp khi có lỗi
			break
		}
	}

	return nil
}

// BroadcastMessages lắng nghe kênh broadcast và gửi tin nhắn tới tất cả các client
func (h *WebSocketHandler) BroadcastMessages() {
	for {
		response := <-h.Broadcast
		for _,u := range response.Clients{
			if clients, ok := h.Clients[u.UserId]; ok {
				for client := range clients {
					err := client.WriteJSON(response)
					if err != nil {
						log.Println("Error writing json:", err)
						client.Close()
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.Clients, u.UserId)
						}
					}
				}
			}
		}
	}
}
