package main

import (
  "fmt"
  "net/http"
  "my-go-project/example"
)


func login() {
  fmt.Println("登录")
}

func userCenter() {
  fmt.Println("个人中心")
}

func logout() {
  fmt.Println("注销")
}

func main() {
  fmt.Println("请输入要执行的操作：")
  fmt.Println(`1：登录
2：个人中心
3：注销`)
  var num int
  fmt.Scan(&num)
  var funcMap = map[int]func(){
    1: login,
    2: userCenter,
    3: logout,
  }
  http.ListenAndServe(":8080", nil)
  funcMap[num]()
	example.SayHello()
}
