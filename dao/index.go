package dao

import (
	"fmt"
	"time"

	"leekbox/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB() {
	dsn := "root:lfl730811@tcp(127.0.0.1:3306)/leekbox?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	} else {
		DB = db
	}
	sqlDB, err := DB.DB()
	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func init() {
	initDB()
	DB.AutoMigrate(&model.User{})
}

func CreateNewUser(user *model.User) (*model.User, error) {
	if err := DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserInfoById(id int) model.User {
	user := model.User{}
	DB.Where("id = ?", id).First(&user)
	return user
}
