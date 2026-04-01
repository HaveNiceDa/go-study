package example

import (
	"fmt"
	// "reflect"
	// "strings"
	"time"
	// "sync"
)
type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	money   int
}

var count = 0

// 缓冲区大小为1的管道
var lock = make(chan struct{}, 1)

func Add() {
    // 加锁
	lock <- struct{}{}
	fmt.Println("当前计数为", count, "执行加法")
	count += 1
    // 解锁
	<-lock
}

func Sub() {
    // 加锁
	lock <- struct{}{}
	fmt.Println("当前计数为", count, "执行减法")
	count -= 1
    // 解锁
	<-lock
}

func write(ch chan<- int) {
	// 只能对管道发送数据
	ch <- 1
}

func Send(ch chan<- int) {
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond)
		ch <- i
	}
}

func SayHello() {
	chA := make(chan int)
	chB := make(chan int)
	chC := make(chan int)
	defer func() {
		close(chA)
		close(chB)
		close(chC)
	}()
	go Send(chA)
	go Send(chB)
	go Send(chC)
    // for循环
	for {
		select {
		case n, ok := <-chA:
			fmt.Println("A", n, ok)
		case n, ok := <-chB:
			fmt.Println("B", n, ok)
		case n, ok := <-chC:
			fmt.Println("B", n, ok)
		}
	}
	 
	}
