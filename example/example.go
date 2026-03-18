package example

import (
	"fmt"
)

func SayHello() {
	s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9} // cap = 9
   s2 := s1[3:4:7]                        // cap = 4 - 3 = 1
   s2 = append(s2, 1)
   fmt.Println(s2)
   fmt.Println(s1)
}
