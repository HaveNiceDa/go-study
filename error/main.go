package main

import (
	"errors"
	"fmt"
)

func div(a, b int) (res int, err error) {
	if b == 0 {
		err = errors.New("division by zero")
		return
	}
	res = a / b
	return
}

func server() (res int, err error) {
	res, err = div(2, 0)
	if err != nil {
		// 把错误向上传递
		return
	}
	// 执行其他的逻辑
	res++
	res *= 2
	return
}

func main() {
	result, err := server()
	if err != nil {
		fmt.Printf("结果: %d, 错误: %v\n", result, err)
	} else {
		fmt.Printf("结果: %d\n", result)
	}
}
