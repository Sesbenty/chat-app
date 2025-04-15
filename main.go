package main

import (
	"chat-app/api/login"
	"chat-app/api/room"
	"chat-app/handlers"
	"chat-app/hub"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	h := hub.NewHub()
	go h.Run()
	router := gin.Default()
	router.LoadHTMLGlob("public/*.html")

	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})

	conten_group := router.Group("/")
	conten_group.Use(login.AuthRedirectMiddleware())
	conten_group.GET("/chat", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	auth_group := router.Group("/api/v1/auth")
	auth_group.POST("/register", login.RegisterHandler)
	auth_group.POST("/login", login.LoginHandler)
	auth_group.GET("/logout", login.LogoutHandler)

	rooms_group := router.Group("/api/v1/rooms")
	rooms_group.Use(login.AuthAPIMiddleware())
	rooms_group.GET("/", room.GetRooms)
	rooms_group.POST("/", room.CreateRoom)
	rooms_group.DELETE("/:id", room.DeleteRoom)

	v1 := router.Group("/api/v1")
	v1.GET("/ws", handlers.WsHandler(h))

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
