package example

import (
	"fmt"
	"sync"
)

func SayHello() {
	var group sync.WaitGroup
	var mutex sync.Mutex
	group.Add(10)
   // map
   mp := make(map[string]int, 10)
   for i := 0; i < 10; i++ {
      go func() {
         // 写操作
         for i := 0; i < 100; i++ {
            mutex.Lock()
            mp["helloworld"] = 1
            mutex.Unlock()
         }
         // 读操作
         for i := 0; i < 10; i++ {
            mutex.Lock()
            fmt.Println(mp["helloworld"])
            mutex.Unlock()
         }
         group.Done()
      }()
   }
   group.Wait()
	
}
