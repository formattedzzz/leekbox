package model

import (
	"time"
)

const (
	UNHANDLED_ERROR     = "发生未知错误"
	USER_EXISTED        = "该用户ID已被占用"
	USER_ACCESSIABLE    = "该用户ID可用"
	USER_SIGNUP_SUCCESS = "用户注册成功"
	USER_LOGIN_SUCCESS  = "用户登录成功"
	USER_NOT_EXISTED    = "用户不存在"
	USER_PASS_INVALID   = "用户密码错误"
	USER_INFO_SUCCEED   = "获取用户信息成功"
)

type User struct {
	Id        int       `json:"id" gorm:"primary_key"`
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
