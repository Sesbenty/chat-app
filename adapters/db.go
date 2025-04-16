package adapters

import (
	"chat-app/models"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DBADDR")
	cfg.DBName = "chat-app"

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
}

func GetRooms() ([]models.Room, error) {
	var rooms []models.Room

	rows, err := db.Query("SELECT * FROM rooms")
	if err != nil {
		return nil, fmt.Errorf("getRooms %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			return nil, fmt.Errorf("getRooms %v", err)
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getRooms %v", err)
	}

	return rooms, nil
}

func GetMessagesByRoom(roomID int64) ([]models.MessageData, error) {
	var messages []models.MessageData
	rows, err := db.Query("SELECT * FROM message WHERE room_id = ?", roomID)
	if err != nil {
		return nil, fmt.Errorf("getMessagesByRoom %d: %v", roomID, err)
	}
	defer rows.Close()
	for rows.Next() {
		var msg models.MessageData
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.RoomID, &msg.Time, &msg.UserID); err != nil {
			return nil, fmt.Errorf("getMessageByRoom %d: %v", roomID, err)
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getMessagesByRoom %d: %v", roomID, err)
	}
	return messages, nil
}

func AddRoom(room models.Room) (int64, error) {
	result, err := db.Exec("INSERT INTO room (name) VALUES(?)", room.Name)
	if err != nil {
		return 0, fmt.Errorf("addRoom: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addRoom: %v", err)
	}

	return id, nil
}

func AddMessage(msg models.MessageData) (int64, error) {
	result, err := db.Exec("INSERT INTO message (room_id, user_id, content, time) VALUES (?, ?, ?, ?)", msg.RoomID, msg.UserID, msg.Content, msg.Time)
	if err != nil {
		return 0, fmt.Errorf("addMessage: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addMessage: %v", err)
	}
	return id, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	row := db.QueryRow("SELECT id, email, password, username FROM user WHERE email = ?", email)
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("getUserByEmail %q no email", email)
		}
		return user, fmt.Errorf("getUserByEmail %q: %v", email, err)
	}

	return user, nil
}

func AddUser(user models.User) (int64, error) {
	result, err := db.Exec("INSERT INTO user (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	if err != nil {
		return 0, fmt.Errorf("addMessage: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addMessage: %v", err)
	}
	return id, nil
}
