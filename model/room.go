package model

import "time"

var ROOM_MSG = struct {
	ROOM_FORBIDDEN string
	ROOM_NOT_EXIST string
	ROOM_NOT_MATCH string
	ROOM_DELETED   string
}{
	ROOM_FORBIDDEN: "你无权修改该讨论组",
	ROOM_NOT_EXIST: "该讨论组不存在",
	ROOM_NOT_MATCH: "讨论组ID与OWNER_ID不匹配",
	ROOM_DELETED:   "该讨论组已被回收",
}

type Room struct {
	Id        int       `json:"id" gorm:"primary_key;autoIncrement"`
	OwnerId   int       `json:"owner_id"`
	Title     string    `json:"title" gorm:"index:idx_room;type:varchar(255);default:''"`
	Desc      string    `json:"desc" gorm:"index:idx_room;type:varchar(255);default:''"`
	Rate      int       `json:"rate" gorm:"type:tinyint;default:0"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(255);default:''"`
	Status    int       `json:"status" gorm:"type:tinyint;default:0"`
	Deleted   int       `json:"deleted" gorm:"type:tinyint;default:0"`
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
