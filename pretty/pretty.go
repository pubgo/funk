package pretty

import (
	"io"
	"strings"
	"sync"

	"github.com/k0kubun/pp/v3"
)

func Println(a ...any) {
	_, _ = pp.Println(a...)
}

func Printf(format string, a ...any) {
	_, _ = pp.Printf(format, a...)
}

func Sprint(a ...any) string {
	return pp.Sprint(a...)
}

func Sprintln(a ...any) string {
	return pp.Sprintln(a...)
}

func Fatal(a ...any) {
	pp.Fatal(a...)
}

func Fatalln(a ...any) {
	pp.Fatalln(a...)
}

func Fatalf(format string, a ...any) {
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

func SimplePrint(v any) string {
	return strings.ReplaceAll(Simple().Sprint(v), "\n", "")
}
