package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data"`
}

func Index(c *gin.Context) {
	c.JSON(200, Response{
		Code: 0,
		Msg:  "success",
		Data: nil,
	})
}

func main() {
	gin.SetMode("release")
	// 初始化路由
	r := gin.Default()
	r.Static("st", "gin/static")
	r.GET("/", Index)
	r.GET("users/:id", func(c *gin.Context) {
  userID := c.Param("id")
  fmt.Println(userID)
})
	r.Run(":8080")
}