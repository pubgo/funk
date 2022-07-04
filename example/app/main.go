package main

import (
	"fmt"
	"github.com/pubgo/funk"
)

// 应用的集成开发, 在最后扑捉panic

func A() string {
	panic("未知错误")
}

func B() string {
	return A()
}

func C() string {
	var a = A()
	if a == "" {
		return B()
	}
	return ""
}

func main() {
	defer funk.RecoverAndExit()

	fmt.Println(C())
}
