package main

import (
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
	r.LoadHTMLGlob("gin/templetes/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.Run(":8080")
}