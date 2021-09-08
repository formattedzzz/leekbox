package model

import (
	"time"
)

var USER_MSG = struct {
	USER_EXISTED      string
	USER_NOT_EXISTED  string
	USER_PASS_INVALID string
	USER_FORBIDDEN    string
}{
	USER_EXISTED:      "该用户ID已被占用",
	USER_NOT_EXISTED:  "用户不存在",
	USER_PASS_INVALID: "用户密码错误",
	USER_FORBIDDEN:    "无权修改该用户",
}

type User struct {
	Id        int       `json:"id" gorm:"primary_key;autoIncrement"`
	Uuid      string    `json:"uuid" gorm:"type:char(36)"`
	Pass      string    `json:"-" gorm:"type:varchar(255) not null;default:''"`
	UserId    string    `json:"name" gorm:"index;type:varchar(255);default:''"`
	NickName  string    `json:"nick_name" gorm:"type:varchar(20);default:''"`
	Desc      string    `json:"desc" gorm:"type:varchar(255);default:''"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(255);default:''"`
	Phone     string    `json:"phone" gorm:"type:varchar(20);default:''"`
	Email     string    `json:"email" gorm:"type:varchar(20);default:''"`
	Rate      float64   `json:"rate" gorm:"type:decimal(5,2);default:0"`
	Balance   int       `json:"balance" gorm:"type:int;default:0"`
	UpdatedAt time.Time `json:"-" gorm:"autoUpdateTime"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
