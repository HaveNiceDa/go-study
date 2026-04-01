package main

import (
	"my-go-project/gin/res/res"
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
	r.GET("/login", func(c *gin.Context) {
		enter.OkWithMsg(c, "登陆成功")
	})
	r.GET("/users", func(c *gin.Context) {
		enter.FailWithCode(c, enter.RoleErrCode)
	})
	r.Run(":8080")
}