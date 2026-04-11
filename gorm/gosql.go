package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	// "log"
)

func goSql() {
  db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/gorm_new_db")
if err != nil {
  fmt.Printf("数据库连接失败 %s", err)
}
err = db.Ping()
if err != nil {
  fmt.Printf("数据库连接失败 %s", err)
}
fmt.Printf("数据库连接成功 %v", db)

}