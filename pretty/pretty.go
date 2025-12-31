package pretty

import (
	"io"
	"strings"
	"sync"

	"github.com/k0kubun/pp/v3"
)

// Println is a wrapper for pp.Println
func Println(a ...any) {
	_, _ = pp.Println(a...)
}

// Printf is a wrapper for pp.Printf
func Printf(format string, a ...any) {
	_, _ = pp.Printf(format, a...)
}

// Sprint is a wrapper for pp.Sprint
func Sprint(a ...any) string {
	return pp.Sprint(a...)
}

// Sprintln is a wrapper for pp.Sprintln
func Sprintln(a ...any) string {
	return pp.Sprintln(a...)
}

// Fatal is a wrapper for pp.Fatal
func Fatal(a ...any) {
	pp.Fatal(a...)
}

// Fatalln is a wrapper for pp.Fatalln
func Fatalln(a ...any) {
	pp.Fatalln(a...)
}

// Fatalf is a wrapper for pp.Fatalf
func Fatalf(format string, a ...any) {
	pp.Fatalf(format, a...)
}

// SetWriter is a wrapper for pp.SetWriter
func SetWriter(o io.Writer) {
	pp.SetDefaultOutput(o)
}

// SetDefaultMaxDepth is a wrapper for pp.SetDefaultMaxDepth
func SetDefaultMaxDepth(v int) {
	pp.SetDefaultMaxDepth(v)
}

// Simple is a simple pretty printer
var Simple = sync.OnceValue(func() *pp.PrettyPrinter {
	printer := pp.New()
	printer.SetColoringEnabled(false)
	printer.SetExportedOnly(false)
	printer.SetOmitEmpty(true)
	printer.SetMaxDepth(3)
	return printer
})

// SimplePrint is a simple pretty printer
func SimplePrint(v any) string {
	return strings.ReplaceAll(Simple().Sprint(v), "\n", "")
}
