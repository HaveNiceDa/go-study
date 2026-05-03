# Elasticsearch 里的索引和 Mapping，到底是什么关系？

刚接触 Elasticsearch 的时候，很多人会被 `索引`、`Mapping`、`文档` 这些概念绕晕。  
如果先用 MySQL 的思路去理解，这几个概念其实没有那么抽象。

先说结论：

- `索引（Index）`：可以理解成 ES 里装数据的地方，类似数据库里的“表”
- `Mapping`：可以理解成这张表的字段定义，类似 MySQL 里的“表结构”
- `文档（Document）`：可以理解成表里的一行数据

也就是说：

- 索引决定“数据放在哪”
- Mapping 决定“数据怎么存、怎么查”

## 1. 先看这段示例代码在做什么

这个示例的入口很简单，先连接 Elasticsearch，然后创建一个索引：

```go
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
	indexs.CreateIndex()
}
```

这段代码的意思很直接：

- `core.EsConnect()`：连接 ES
- `indexs.CreateIndex()`：创建一个叫 `user_index` 的索引

对应的创建索引逻辑如下：

```go
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
```

这里最关键的一句是：

```go
BodyString(models.UserModel{}.Mapping())
```

它的作用就是：把 `UserModel` 里定义好的 Mapping 规则交给 Elasticsearch，让 ES 按照这些字段规则来创建索引。

## 2. 索引到底是什么

索引可以先粗暴理解成一张“表”。

比如这里创建的索引叫：

```go
CreateIndex("user_index")
```

那么这个 `user_index` 就可以理解成“用户数据表”。  
以后和用户相关的数据，都会往这个索引里写。

比如一条用户数据可能长这样：

```json
{
  "id": 1,
  "user_name": "zhangsan",
  "nick_name": "张三",
  "created_at": "2026-05-01 12:00:00"
}
```

这条数据在 ES 里就是一个 `document`，它会被保存到 `user_index` 这个索引中。

所以你可以把它们的关系理解成：

- `user_index` 是容器
- 用户数据是一条条放进去的内容

## 3. Mapping 到底是什么

如果说索引是“表”，那 Mapping 就是“建表时的字段规则”。

在你的代码里，`UserModel` 不只是一个 Go 结构体，它还顺带定义了这个索引的 Mapping：

```go
package models

import "time"

type UserModel struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"user_name"`
	NickName  string    `json:"nick_name"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserModel) Mapping() string {
	return `
		{
  "mappings": {
    "properties": {
      "nick_name": { 
        "type": "text"
      },
      "user_name": { 
        "type": "keyword"
      },
      "id": {
        "type": "integer"
      },
      "created_at":{
        "type": "date",
        "null_value": "null",
        "format": "[yyyy-MM-dd HH:mm:ss]"
      }
    }
  }
}
	`
}
```

这个 Mapping 主要告诉 ES 4 件事：

- `id` 是整数
- `user_name` 是精确匹配字段
- `nick_name` 是可搜索、可分词的字段
- `created_at` 是日期字段

所以 Mapping 的作用不是简单的“声明字段”，而是告诉 ES：

- 这个字段是什么类型
- 这个字段该怎么建立索引
- 这个字段更适合怎么搜索

## 4. 为什么字段类型这么重要

Elasticsearch 和传统数据库最大的不同点，就是它特别强调“搜索能力”。

也正因为这样，字段类型会直接影响查询效果。

### `nick_name` 为什么是 `text`

```json
"nick_name": {
  "type": "text"
}
```

`text` 类型适合做全文搜索，特点是会分词。

比如数据是：

```json
{
  "nick_name": "张三同学"
}
```

那么搜索时：

- 搜 `张三同学`，通常可以查到
- 搜 `张三`，也有机会查到
- 搜 `同学`，也可能查到

因为它更偏向“搜内容”。

### `user_name` 为什么是 `keyword`

```json
"user_name": {
  "type": "keyword"
}
```

`keyword` 适合精确匹配，不分词。

比如：

```json
{
  "user_name": "zhangsan"
}
```

那它更适合这种场景：

- 按用户名完整查找
- 判断某个用户名是否存在
- 做聚合、排序、过滤

所以：

- `text` 更像“模糊搜”
- `keyword` 更像“精确查”

## 5. 索引和 Mapping 的关系

这两个概念最好不要分开记，要连起来理解。

可以把它们看成：

- 索引是房子
- Mapping 是房子的户型图

或者用开发更熟悉的话说：

- 索引是存数据的容器
- Mapping 是这个容器的字段规则

也就是说，创建索引的时候，最好同时把 Mapping 设计好。  
这样 Elasticsearch 才知道后面进来的数据应该按什么方式去存、去建倒排索引、去支持搜索。

如果只创建索引，不定义 Mapping，ES 会尝试自动推断字段类型。  
虽然这样很省事，但实际开发里经常会带来问题，比如：

- 本来想精确匹配，结果被当成全文检索
- 日期格式识别和预期不一致
- 字段类型推断错误，后续查询行为变得奇怪

所以更稳妥的方式是：

1. 先设计索引
2. 再写好 Mapping
3. 最后再写入数据

## 6. 用 MySQL 思维理解最简单

如果你刚学 ES，可以先直接这样类比：

- `Index` 类似一张表
- `Mapping` 类似 `CREATE TABLE` 时的字段定义
- `Document` 类似表里的一行记录

比如可以类比成下面这种感觉：

```sql
CREATE TABLE user_index (
  id INT,
  user_name VARCHAR(255),
  nick_name TEXT,
  created_at DATETIME
);
```

当然，Elasticsearch 不是关系型数据库，它更偏向搜索引擎。  
但在入门阶段，这种理解方式是最容易建立直觉的。

## 7. 结合这份代码，整个流程是怎样串起来的

把这几个文件连起来看，流程就很清楚了：

### 第一步：连接 Elasticsearch

在 `main.go` 中调用：

```go
core.EsConnect()
```

这一步负责把 `global.ESClient` 初始化好。

### 第二步：创建索引

继续调用：

```go
indexs.CreateIndex()
```

这一步会去创建 `user_index`。

### 第三步：把 Mapping 传给 ES

在 `CreateIndex()` 里：

```go
BodyString(models.UserModel{}.Mapping())
```

这里把 `UserModel` 对应的 Mapping JSON 发给 Elasticsearch。

### 第四步：ES 按规则生成索引结构

接下来 ES 就知道：

- `id` 应该按整数处理
- `user_name` 应该按精确匹配处理
- `nick_name` 应该按分词搜索处理
- `created_at` 应该按日期处理

这样你后面再写入用户数据时，搜索行为就有了明确规则。

## 8. 一句话总结

一句话记忆：

- `索引` 是“数据放哪里”
- `Mapping` 是“数据怎么存、怎么查”

所以它们不是并列关系，而是“容器”和“规则”的关系。

Elasticsearch 真正厉害的地方，不只是能存数据，而是它能根据 Mapping 的设计，把普通数据变成“可搜索的数据”。

## 9. 实战补充

上面的思路是对的，但如果你直接把 Mapping 当成 JSON 发送给 ES，要注意一个细节：

- JSON 本身不能写 `//` 注释

也就是说，如果是实际发给 Elasticsearch 的请求体，Mapping 字符串里不要保留这类注释，否则会导致索引创建失败。

如果只是博客展示、学习理解，写注释没有问题；如果是程序真正发送给 ES 的内容，就要保证它是合法 JSON。
