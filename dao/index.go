package dao

import (
	"fmt"
	"time"

	"leekbox/config"
	"leekbox/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDB struct {
	DB *gorm.DB
}

func New(conf config.Configuration) (*GormDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/leekbox?charset=utf8mb4&parseTime=True&loc=Local", conf.DB_USER, conf.DB_PASS, conf.DB_HOST, conf.DB_PORT)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.User{}, &model.Room{}, &model.Comment{}); err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return &GormDB{db}, nil
}
