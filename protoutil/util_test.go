package gen

import (
	"fmt"
	"strings"
	"testing"

	"github.com/iancoleman/strcase"
)

func importHandle(pkg string) string {
	if strings.Contains(pkg, "/") {
		names := strings.Split(pkg, "/")
		pkg = names[0]
		for _, name := range names[1:] {
			pkg += strcase.ToCamel(name)
		}
	}

	if strings.Contains(pkg, ".") {
		names := strings.Split(pkg, ".")
		pkg = names[0]
		for _, name := range names[1:] {
			pkg += strcase.ToCamel(name)
		}
	}

	return pkg
}

func TestImportHandle(t *testing.T) {
	fmt.Println(importHandle("hello.v1"))
	fmt.Println(importHandle("hello.v1.v2"))
	fmt.Println(importHandle("hello/v1/v2"))
}
