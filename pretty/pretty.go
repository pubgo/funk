package pretty

import (
	"io"
	"strings"
	"sync"

	"github.com/k0kubun/pp/v3"
)

func Println(a ...interface{}) {
	_, _ = pp.Println(a...)
}

func Printf(format string, a ...interface{}) {
	_, _ = pp.Printf(format, a...)
}

func Sprint(a ...interface{}) string {
	return pp.Sprint(a...)
}

func Sprintln(a ...interface{}) string {
	return pp.Sprintln(a...)
}

func Fatal(a ...interface{}) {
	pp.Fatal(a...)
}

func Fatalln(a ...interface{}) {
	pp.Fatalln(a...)
}

func Fatalf(format string, a ...interface{}) {
	pp.Fatalf(format, a...)
}

func SetWriter(o io.Writer) {
	pp.SetDefaultOutput(o)
}

func SetDefaultMaxDepth(v int) {
	pp.SetDefaultMaxDepth(v)
}

var Simple = sync.OnceValue(func() *pp.PrettyPrinter {
	printer := pp.New()
	printer.SetColoringEnabled(false)
	printer.SetExportedOnly(false)
	printer.SetOmitEmpty(true)
	printer.SetMaxDepth(3)
	return printer
})

func SimplePrint(v interface{}) string {
	return strings.ReplaceAll(Simple().Sprint(v), "\n", "")
}
