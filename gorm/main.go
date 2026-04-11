package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
  ID   uint
  Name string
}

func main() {
  // 连接数据库
  db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/gorm_new_db"), &gorm.Config{})
  if err != nil {
    fmt.Printf("数据库连接失败 %s", err)
  }
  fmt.Printf("数据库连接成功 %v", db)

	var userList []User
  db.Find(&userList)

  fmt.Println(userList)
}