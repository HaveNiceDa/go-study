package indexs

import (
	"context"
	"fmt"

	"my-go-project/es/global"
	"my-go-project/es/models"
)

func CreateIndex() {
	createIndex, err := global.ESClient.
		CreateIndex("user_index").
		BodyString(models.UserModel{}.Mapping()).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(createIndex)
	fmt.Println("索引创建成功")
}
