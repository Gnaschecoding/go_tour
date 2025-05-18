package main

import (
	"fmt"
	"unsafe"
)

func main() {
	v := 1
	fmt.Printf("类型: %T\n", v)                   // 输出类型
	fmt.Printf("占用字节数: %d\n", unsafe.Sizeof(v)) // 输出占用字节数
}
