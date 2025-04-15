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
	router.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", nil)
	})

	content_group := router.Group("/")
	content_group.Use(login.AuthRedirectMiddleware())
	content_group.GET("/chat", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	content_group.GET("/room", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "rooms.html", nil)
	})
	content_group.GET("/room/:id", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "room.html", nil)
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
	rooms_group.PUT("/:id", room.UpdateRoom)
	rooms_group.GET("/:id/messages", room.GetMessages)

	v1 := router.Group("/api/v1")
	v1.Use(login.AuthAPIMiddleware())
	v1.GET("/ws", handlers.WsHandler(h))

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
