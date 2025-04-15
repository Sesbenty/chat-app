package handlers

import (
	"chat-app/client"
	"chat-app/hub"
	"chat-app/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development. In production, restrict this!
		return true
	},
}

// WsHandler handles WebSocket connections.
func WsHandler(h *hub.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		new_client := &hub.Client{Conn: conn, Send: make(chan models.Message, 256)}
		if err != nil {
			log.Println(err)
			return
		}
		h.Register <- new_client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump(h, new_client)
		go client.ReadPump(h, new_client)
	}
}
