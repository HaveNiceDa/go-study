package main

import (
	"fmt"
	"sync"
)

var (
	count int
	rw    sync.RWMutex // 读写锁
)

// 读操作：可以并发
func read() {
	rw.RLock() // 读锁
	defer rw.RUnlock()

	// 读数据
	_ = count
	// 模拟读取数据
	fmt.Printf("读取到数据: %d\n", count)
	// time.Sleep(10 * time.Millisecond)
}

// 写操作：独占
func write() {
	rw.Lock() // 写锁
	defer rw.Unlock()

	count++ // 修改数据
	fmt.Printf("写入数据: %d\n", count)
	// time.Sleep(100 * time.Millisecond)
}

func main() {
	var wg sync.WaitGroup

	// 启动 100 个读协程（并发执行，超快）

	go func() {
		for i := 0; i < 100; i++ {
		wg.Add(1)
		if i == 0 {
				// time.Sleep(time.Millisecond * 1000)
			}
		go func() {
			
			defer wg.Done()
			read()
		}()
	}
	}()
	

	// 启动 10 个写协程（串行，独占）
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			write()
		}()
	}

	wg.Wait()
}