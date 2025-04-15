package room

import (
	"chat-app/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomRequest struct {
	Name string `json:"name" binding:"required"`
}

type RoomResponse struct {
	Message string `json:"message"`
}

var RoomDatabase = map[int64]models.Room{
	1: {ID: 1, Name: "Room 1"},
	2: {ID: 2, Name: "Room 2"},
}

var id int64 = 3

func GetRooms(c *gin.Context) {
	c.JSON(200, RoomDatabase)
}

func CreateRoom(c *gin.Context) {
	var request RoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	room := models.Room{
		ID:   id,
		Name: request.Name,
	}
	id++
	RoomDatabase[room.ID] = room
	c.JSON(200, request)
}

func DeleteRoom(c *gin.Context) {
	roomID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	delete(RoomDatabase, roomID)
	c.JSON(200, gin.H{"message": "Room deleted"})
}

func UpdateRoom(c *gin.Context) {
	roomID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	RoomDatabase[roomID] = room
	c.JSON(200, room)
}
