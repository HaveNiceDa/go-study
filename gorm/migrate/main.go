package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserModelss struct {
  ID        int64     `gorm:"primaryKey"`      // 主键
  Name      string    `gorm:"not null;unique"` // 不能为空，且唯一
  CreatedAt time.Time // 在创建记录时自动设置为当前时间。
}

func main() {
  // 连接数据库
  db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/gorm_new_db"), &gorm.Config{})
  if err != nil {
    fmt.Println(err)
    return
  }
  
  err = db.AutoMigrate(&UserModelss{})
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println("生成表结构成功！")
}