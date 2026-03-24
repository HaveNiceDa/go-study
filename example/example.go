package example

import (
	"fmt"
	"os"
)
// 起重机接口
type Crane interface { 
	JackUp() string
	Hoist() string
}

// 起重机A
type CraneA struct {
	work int //内部的字段不同代表内部细节不一样
}

func (c CraneA) Work() {
	fmt.Println("使用技术A")
}
func (c CraneA) JackUp() string {
	c.Work()
	return "jackup"
}

func (c CraneA) Hoist() string {
	c.Work()
	return "hoist"
}

// 起重机B
type CraneB struct {
	boot string
}

func (c CraneB) Boot() {
	fmt.Println("使用技术B")
}

func (c CraneB) JackUp() string {
	c.Boot()
	return "jackup"
}

func (c CraneB) Hoist() string {
	c.Boot()
	return "hoist"
}

type ConstructionCompany struct {
	Crane Crane // 只根据Crane类型来存放起重机
}

func (c *ConstructionCompany) Build() {
	fmt.Println(c.Crane.JackUp())
	fmt.Println(c.Crane.Hoist())
	fmt.Println("建筑完成")
}
type Person struct {
    Name string  // 导出字段
    age  int     // 非导出字段
}

type GenericSlice[T int | int32 | int64] []T

type SignedInt interface {
   int8 | int16 | int | int32 | int64
}

type Integer interface {
   int8 | int16 | int | int32 | int64 | uint8 | uint16 | uint | uint32 | uint64
}

type Number interface {
	SignedInt
	Integer
}
func Do[T Number](n T) T {
   return n
}
func dangerOp(str string) {
	if len(str) == 0 {
		fmt.Println("fatal")
		defer os.Exit(1)
	}
	fmt.Println("正常逻辑")
}
func SayHello() {
	dangerOp("")
	Do[int](2)

	// 创建一个整数切片
	intSlice := GenericSlice[int]{1, 2, 3}
	fmt.Println("整数切片:", intSlice)

	// 使用起重机A
	company := ConstructionCompany{CraneA{}}
	company.Build()
	fmt.Println()
	// 更换起重机B
	company.Crane = CraneB{}
	company.Build()
	fmt.Println()
	
	// 演示值接收者和指针接收者的区别
	fmt.Println("=== 演示值接收者和指针接收者的区别 ===")
	
	// 创建公司实例
	company2 := ConstructionCompany{CraneA{}}
	fmt.Println("初始公司:", company2.Crane)
	
	// 使用值接收者方法（不会修改原始值）
	company2.UpdateWithValue("Value Receiver")
	fmt.Println("值接收者修改后:", company2.Crane)
	
	// 使用指针接收者方法（会修改原始值）
	company2.UpdateWithPointer("Pointer Receiver")
	fmt.Println("指针接收者修改后:", company2.Crane)
	
	// 演示导出字段和非导出字段的打印区别
	fmt.Println("\n=== 演示导出字段和非导出字段的打印区别 ===")
	
	// 包含导出和非导出字段的结构体
	p := Person{Name: "Tom", age: 20}
	fmt.Println("Person结构体（混合字段）:", p)
	
	// 对比只有非导出字段的情况
	type NonExported struct {
		a int
		b string
	}
	ne := NonExported{a: 1, b: "test"}
	fmt.Println("只有非导出字段:", ne)
	
	// 对比只有导出字段的情况
	type Exported struct {
		A int
		B string
	}
	e := Exported{A: 1, B: "test"}
	fmt.Println("只有导出字段:", e)
}

// 值接收者方法（不会修改原始值）
func (c ConstructionCompany) UpdateWithValue(name string) {
	// 尝试修改字段（只会修改副本）
	c.Crane = CraneB{boot: name}
	fmt.Println("值接收者内部:", c.Crane)
}

// 指针接收者方法（会修改原始值）
func (c *ConstructionCompany) UpdateWithPointer(name string) {
	// 修改字段（会修改原始值）
	c.Crane = CraneB{boot: name}
	fmt.Println("指针接收者内部:", c.Crane)
}
