package indexs

import (
	"context"
	"fmt"

	"my-go-project/es/global"
)

// ExistsIndex 判断索引是否存在
func ExistsIndex(index string) bool {
	exists, err := global.ESClient.
		IndexExists(index).
		Do(context.Background())
	fmt.Println(exists, err)
	return exists
}
