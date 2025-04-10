package main

import (
	"chat-app/handlers"
	"chat-app/hub"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new hub.
	h := hub.NewHub()
	go h.Run()

	// Create a new Gin router.
	router := gin.Default()

	// Serve static files (optional, for a frontend).
	router.StaticFile("/", "public/index.html")

	// WebSocket route.
	router.GET("/ws", handlers.WsHandler(h))

	// Start the server.
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
