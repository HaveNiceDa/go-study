# Go 学习项目

这是一个用于学习 Go 语言基础的项目，展示了 Go 的基本项目结构和包管理。

## 项目结构

```
my-go-project/
├── go.mod              # Go 模块定义文件
├── main/
│   └── main.go         # 主程序入口
└── example/
    └── example.go      # 示例包
```

## 快速开始

### 前置要求

- Go 1.26.1 或更高版本

### 运行项目

```bash
# 运行主程序
go run main/main.go
```

输出：
```
Hello
```

### 编译项目

```bash
# 编译为可执行文件
go build -o my-app main/main.go

# 运行编译后的程序
./my-app
```

## 项目说明

### 包结构

- **main**：主程序包，包含程序的入口函数 `main()`
- **example**：示例包，展示了如何创建和使用自定义包

### Go 模块

项目使用 Go 模块进行依赖管理，模块名为 `my-go-project`。

## Go 基础知识点

### 1. 包（Package）

- 每个 Go 文件必须声明所属的包
- 同一目录下的所有文件必须属于同一个包
- `package main` 是特殊的包，包含程序的入口点

### 2. 导入（Import）

使用 `import` 关键字导入其他包：
```go
import "fmt"
import "my-go-project/example"
```

### 3. 函数定义

```go
func functionName() {
    // 函数体
}
```

### 4. 格式化输出

使用 `fmt` 包进行输出：
```go
fmt.Println("Hello")  // 打印并换行
fmt.Print("Hello")    // 打印不换行
```

## 学习资源

- [Go 官方文档](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)

## 许可证

MIT License
