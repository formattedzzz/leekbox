package main

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/davecgh/go-spew/spew"
)

// 这里的router本质就是一个struct类型为gin.Engine的指针
func initRouter(router *gin.Engine) {
	router.LoadHTMLFiles("templates/index.html")
}
func leo12() {
	s1 := []string{
		"leo",
		"bob",
		"npm",
	}
	m1 := map[string]interface{}{
		"name": "leooo",
		"age":  2,
	}
	// spew 包是更好的输出映射类型的值
	spew.Dump(s1, m1)
	spew.Printf("v: %v\n", m1)
}
func main() {
	leo12()
	// 在必要的模块 panic之后需要重启一下
	// defer func() {
	// 	if fatal := recover(); fatal != nil {
	// 		fmt.Println("panic captured", fatal)
	// 		return
	// 	}
	// }()

	router := gin.Default()
	router.Static("/static", "./static")
	// router.LoadHTMLGlob("./templates/**/*")
	router.SetFuncMap(template.FuncMap{
		"format": func(str string) string {
			return "leo-" + str
		},
	})
	initRouter(router)
	router.GET("/", func(ctx *gin.Context) {
		// ctx.JSON(200, map[string]interface{}{
		// 	"code": 20000,
		// 	"data": []interface{}{
		// 		1, "", true, 1.3,
		// 	},
		// 	"msg": "success",
		// })
		ctx.DefaultQuery("userid", "0")
		name := ctx.Query("name")
		type info struct {
			code int
		}
		ctx.HTML(200, "index.html", map[string]interface{}{
			"data":   []interface{}{1, 2, 3, 4, 5, info{1}},
			"msg":    "success",
			"name":   name,
			"userid": ctx.Query("userid"),
		})
	})
	router.GET("/user/info/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Println("id not found")
		}
		c.JSON(200, gin.H{
			"code":    20000,
			"messgae": "success",
			"data":    []string{"leo", "bob"},
			"id":      id,
		})
	})
	router.Run(":7003")
}
