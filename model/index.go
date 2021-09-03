package model

import "time"

type Resp struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Return(code int, data interface{}, message string) Resp {
	return Resp{code, data, message}
}

const (
	API_SUCCESS     = "请求成功"
	UNHANDLED_ERROR = "发生未知错误"
)

type Comment struct {
	Id        int       `json:"id" gorm:"primary_key;autoIncrement"`
	Uid       int       `json:"uid" gorm:"not null"`
	RoomId    int       `json:"room_id" gorm:"not null"`
	Type      int       `json:"type" gorm:"type:tinyint;default:0"`
	Content   string    `json:"content" gorm:"type:text"`
	Attach    string    `json:"attach" gorm:"type:text"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
