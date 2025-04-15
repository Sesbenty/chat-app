package models

type Room struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

const (
	Admin  = iota
	Member = iota
)

type UserRoomMember struct {
	RoomID int64 `json:"room_id"`
	UserID int64 `json:"user_id"`
	Role   int   `json:"role"`
}

type MessageData struct {
	ID      int64  `json:"id"`
	RoomID  int64  `json:"room_id"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
	Time    string `json:"time"`
}
