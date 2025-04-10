package client

import (
	"chat-app/hub"
	"chat-app/models"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

func ReadPump(h *hub.Hub, c *hub.Client) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var msg models.Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}
		msg.Timestamp = time.Now()
		h.Broadcast <- msg
	}
}

func WritePump(h *hub.Hub, c *hub.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			messageJson, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshalling message: %v", err)
				return
			}
			w.Write(messageJson)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				messageJson, err := json.Marshal(<-c.Send)
				if err != nil {
					log.Printf("error marshalling message: %v", err)
					return
				}
				w.Write(messageJson)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
