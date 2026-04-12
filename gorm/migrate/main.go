package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserModel struct {
  ID        int64     `gorm:"primaryKey"`      // 主键
  Name      string    `gorm:"not null;unique"` // 不能为空，且唯一
  CreatedAt time.Time // 在创建记录时自动设置为当前时间。
}

// InsertData 插入数据
func InsertData(db *gorm.DB) {
  user := UserModel{Name: "张三"}
  err := db.Create(&user).Error
  // 创建成功之后，数据会回填
  fmt.Println(user, err)
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
  fmt.Println("创建的钩子函数")
  u.Name = "枫枫"
  return nil
}

func main() {
  // 连接数据库
  db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/gorm_new_db?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
  if err != nil {
    fmt.Println(err)
    return
  }
  
  // 查全部
var users []UserModel
db.Debug().Find(&users)
fmt.Println(users)
}