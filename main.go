package main

import (
	"fmt"
	"leekbox/config"
	"leekbox/dao"
	_ "leekbox/dao"
	"leekbox/model"
	"leekbox/router"
	"time"
)

var LOC *time.Location

func init() {
	initLoc()
}

func initLoc() {
	if lac, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		LOC = lac
	}
}

// @title LeekBox API
// @version 1.0
// @description Leekbox. a fabulous share-room.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	recover := func() {
		// 在必要的模块 panic之后需要重启一下
		if fatal := recover(); fatal != nil {
			fmt.Println("panic captured", fatal)
			return
		}
	}
	defer recover()

	config := config.Get()
	tableList := []interface{}{&model.User{}, &model.Room{}, &model.Comment{}}
	db, err_db := dao.New(*config, tableList)
	if err_db != nil {
		panic(fmt.Errorf("数据库初始化失败%s", err_db))
	}
	app := router.Create(db, *config)
	fmt.Println("----------LEEK_BOX----------")
	fmt.Printf("-%26s-\n", "")
	fmt.Printf("-%26s-\n", "")
	fmt.Printf("-%26s-\n", "")
	fmt.Printf("-%26s-\n", "")
	fmt.Printf("-%26s-\n", "")
	fmt.Printf("------------%5s-----------\n", config.VERSION)
	err := app.Run(fmt.Sprintf(":%d", config.PORT))
	if err != nil {
		fmt.Printf("go-app启动失败 %s", err)
	}
}
