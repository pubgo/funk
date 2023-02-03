package recovery

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/alecthomas/repr"
	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/generic"
	"github.com/pubgo/funk/log"
	"github.com/pubgo/funk/result"
)

func Result[T any](ret *result.Result[T]) {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	*ret = result.Err[T](err)
}

func Err(gErr *error, fns ...func(err error) error) {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	if len(fns) > 0 && fns[0] != nil {
		*gErr = fns[0](err)
		return
	}

	*gErr = err
}

func Raise(fns ...func(err error) error) {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	if len(fns) > 0 && fns[0] != nil {
		panic(errors.WrapCaller(fns[0](err), 1))
	}

	panic(errors.WrapCaller(err, 1))
}

func Recovery(fn func(err error)) {
	assert.If(fn == nil, "[fn] should not be nil")

	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	fn(err)
}

func Exit(handlers ...func()) {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	if len(handlers) > 0 {
		handlers[0]()
	}

	errors.Debug(err)
	debug.PrintStack()
	os.Exit(1)
}

func DebugPrint() {
	err := errors.Parse(recover())
	if generic.IsNil(err) {
		return
	}

	errors.Debug(err)
	debug.PrintStack()
}

func Dump() {
	var errS string
	if err := errors.Parse(recover()); !generic.IsNil(err) {
		var bytes, err11 = json.MarshalIndent(err, "  ", "  ")
		if err11 != nil {
			errS = repr.String(err11)
		} else {
			errS = string(bytes)
		}
	}

	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, true)

	buf = []byte(strings.ReplaceAll("  "+string(buf[:n]), "\n", "\n\t"))
	buf = []byte(strings.ReplaceAll(string(buf), "\tgoroutine ", "  goroutine "))
	path := os.Getenv("SIGDUMP_PATH")
	if path == "" {
		path = fmt.Sprintf("/tmp/sigdump-%d.log", os.Getpid())
	}

	log.Info().Msgf("dump path:%s", path)

	var w *os.File
	if path == "-" {
		w = os.Stdout
	} else if path == "+" {
		w = os.Stderr
	} else {
		var err error
		w, err = os.Create(path)
		if err != nil {
			return
		}

		defer func(w *os.File) {
			err = w.Close()
			if err != nil {
				return
			}
		}(w)
	}

	now := time.Now().Format(time.RFC3339)
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	ppid := os.Getppid()
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("Sigdump time=%s host=%s pid=%d ppid=%d \n%s\n", now, hostname, pid, ppid, buf[:n]))

	sb.WriteString("\n  Mem Stat:\n")
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	sb.WriteString(fmt.Sprintf("\tAlloc = %v\n", memStats.Alloc))
	sb.WriteString(fmt.Sprintf("\tTotalAlloc = %v\n", memStats.TotalAlloc))
	sb.WriteString(fmt.Sprintf("\tSys = %v\n", memStats.Sys))
	sb.WriteString(fmt.Sprintf("\tLookups = %v\n", memStats.Lookups))
	sb.WriteString(fmt.Sprintf("\tMallocs = %v\n", memStats.Mallocs))
	sb.WriteString(fmt.Sprintf("\tFrees = %v\n", memStats.Frees))
	sb.WriteString(fmt.Sprintf("\tHeapAlloc = %v\n", memStats.HeapAlloc))
	sb.WriteString(fmt.Sprintf("\tHeapSys = %v\n", memStats.HeapSys))
	sb.WriteString(fmt.Sprintf("\tHeapIdle = %v\n", memStats.HeapIdle))
	sb.WriteString(fmt.Sprintf("\tHeapInuse = %v\n", memStats.HeapInuse))
	sb.WriteString(fmt.Sprintf("\tHeapReleased = %v\n", memStats.HeapReleased))
	sb.WriteString(fmt.Sprintf("\tHeapObjects = %v\n", memStats.HeapObjects))
	sb.WriteString(fmt.Sprintf("\tStackInuse = %v\n", memStats.StackInuse))
	sb.WriteString(fmt.Sprintf("\tStackSys = %v\n", memStats.StackSys))
	sb.WriteString(fmt.Sprintf("\tMSpanInuse = %v\n", memStats.MSpanInuse))
	sb.WriteString(fmt.Sprintf("\tMSpanSys = %v\n", memStats.MSpanSys))
	sb.WriteString(fmt.Sprintf("\tMCacheInuse = %v\n", memStats.MCacheInuse))
	sb.WriteString(fmt.Sprintf("\tMCacheSys = %v\n", memStats.MCacheSys))
	sb.WriteString(fmt.Sprintf("\tBuckHashSys = %v\n", memStats.BuckHashSys))
	sb.WriteString(fmt.Sprintf("\tGCSys = %v\n", memStats.GCSys))
	sb.WriteString(fmt.Sprintf("\tOtherSys = %v\n", memStats.OtherSys))
	sb.WriteString(fmt.Sprintf("\tNextGC = %v\n", memStats.NextGC))
	sb.WriteString(fmt.Sprintf("\tLastGC = %v\n", memStats.LastGC))
	sb.WriteString(fmt.Sprintf("\tPauseTotalNs = %v\n", memStats.PauseTotalNs))
	sb.WriteString(fmt.Sprintf("\tNumGC = %v\n", memStats.NumGC))
	sb.WriteString(fmt.Sprintf("\tGCCPUFraction = %v\n", memStats.GCCPUFraction))
	sb.WriteString(fmt.Sprintf("\tDebugGC = %v\n", memStats.DebugGC))
	sb.WriteString(fmt.Sprintf("\t%s\t", errS))
	_, err := w.WriteString(sb.String())
	if err != nil {
		return
	}
}
