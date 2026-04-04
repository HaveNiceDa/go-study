package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("协程 %d 开始\n", id)
	time.Sleep(time.Second)
	fmt.Printf("协程 %d 结束\n", id)
}

func main() {
	var wg sync.WaitGroup

	// 初始启动2个协程
	fmt.Println("启动初始协程...")
	wg.Add(2)
	go worker(1, &wg)
	go worker(2, &wg)

	// 模拟运行中发现需要更多协程
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("发现需要更多协程，动态增加...")
		wg.Add(2) // 动态增加2个协程
		go worker(3, &wg)
		go worker(4, &wg)
	}()

	wg.Wait()
	fmt.Println("全部完成")
}