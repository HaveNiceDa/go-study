package main

import (
	"fmt"
	"my-go-project/es/core"
	"my-go-project/es/global"
)

func main() {
	core.EsConnect()
	fmt.Println(global.ESClient)
}
