package model

import (
	"time"
)

var USER_MSG = struct {
	USER_EXISTED      string
	USER_NOT_EXISTED  string
	USER_PASS_INVALID string
}{
	USER_EXISTED:      "该用户ID已被占用",
	USER_NOT_EXISTED:  "用户不存在",
	USER_PASS_INVALID: "用户密码错误",
}

type User struct {
	Id        int       `json:"id" gorm:"primary_key;autoIncrement"`
	Uuid      string    `json:"uuid" gorm:"type:char(36)"`
	Pass      string    `json:"omit" gorm:"type:varchar(255) not null;default:''"`
	UserId    string    `json:"name" gorm:"index;type:varchar(255);default:''"`
	NickName  string    `json:"nick_name" gorm:"type:varchar(20);default:''"`
	Desc      string    `json:"desc" gorm:"type:varchar(255);default:''"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(255);default:''"`
	Phone     string    `json:"phone" gorm:"type:varchar(20);default:''"`
	Rate      float64   `json:"rate" gorm:"type:decimal(5,2);default:0"`
	Balance   int       `json:"balance" gorm:"type:int;default:0"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type Room struct {
	Id        int       `json:"id" gorm:"primary_key;autoIncrement"`
	OwnerId   int       `json:"owner_id"`
	Title     string    `json:"name" gorm:"index:idx_room;type:varchar(255);default:''"`
	Desc      string    `json:"desc" gorm:"index:idx_room;type:varchar(255);default:''"`
	Rate      int       `json:"rate" gorm:"type:tinyint;default:0"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(255);default:''"`
	Status    int       `json:"status" gorm:"type:tinyint;default:0"`
	Deleted   bool      `json:"deleted" gorm:"type:tinyint;default:0"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type RoomInfo struct {
	Room
	IsOwner bool `json:"is_owner"`
	Owner   User `gorm:"foreignKey:OwnerId;references:Id" json:"owner,omitempty"`
}

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
