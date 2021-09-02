package dao

import (
	"fmt"
	"time"

	"leekbox/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDB struct {
	DB *gorm.DB
}

func New(conf config.Configuration, tableList []interface{}) (*GormDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/leekbox?charset=utf8mb4&parseTime=True&loc=Local", conf.DB_USER, conf.DB_PASS, conf.DB_HOST, conf.DB_PORT)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(tableList...); err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return &GormDB{db}, nil
}
