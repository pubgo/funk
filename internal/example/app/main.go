package main

import (
	"fmt"

	"github.com/pubgo/funk/recovery"
)

// 应用的集成开发, 在最后扑捉panic

func A() string {
	panic("未知错误")
}

func B() string {
	return A()
}

func C() string {
	a := A()
	if a == "" {
		return B()
	}
	return ""
}

func main() {
	defer recovery.Exit()

	fmt.Println(C())
}
