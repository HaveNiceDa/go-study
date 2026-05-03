package main

import (
	"fmt"
	"my-go-project/es/core"
	"my-go-project/es/global"
	"my-go-project/es/indexs"
)

func main() {
	core.EsConnect()
	fmt.Println(global.ESClient)
	indexs.ExistsIndex("user_index")
	indexs.CreateIndex()

}
