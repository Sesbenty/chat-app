package room

import (
	"chat-app/models"
	"maps"
	"slices"
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

var MessagesDatabase = map[int64]models.MessageData{
	1:  {ID: 1, RoomID: 1, Content: "Hello from Room 1", UserID: 1},
	2:  {ID: 2, RoomID: 1, Content: "Hi there!", UserID: 2},
	3:  {ID: 3, RoomID: 1, Content: "How are you?", UserID: 1},
	4:  {ID: 4, RoomID: 2, Content: "Welcome to Room 2", UserID: 3},
	5:  {ID: 5, RoomID: 2, Content: "Nice to be here", UserID: 4},
	6:  {ID: 6, RoomID: 2, Content: "What's up?", UserID: 3},
	7:  {ID: 7, RoomID: 1, Content: "Goodbye!", UserID: 2},
	8:  {ID: 8, RoomID: 1, Content: "See you later", UserID: 1},
	9:  {ID: 9, RoomID: 2, Content: "Bye everyone", UserID: 4},
	10: {ID: 10, RoomID: 2, Content: "Have a good day!", UserID: 3},
}

var id int64 = 3

func GetRooms(c *gin.Context) {
	c.JSON(200, slices.Collect(maps.Values(RoomDatabase)))
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

func GetMessages(c *gin.Context) {
	roomID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	messages := []models.MessageData{}
	for k, value := range MessagesDatabase {
		if value.RoomID == roomID {
			messages = append(messages, MessagesDatabase[k])
		}
	}
	c.JSON(200, messages)	
	}
