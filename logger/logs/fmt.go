package logs

import (
	"bytes"
	"fmt"
	"github.com/go-logfmt/logfmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pubgo/funk/logger"
)

// NewTextLog returns a deferred logger that writes in logfmt using the default options.
// The logger defers configuring its options until it is instantiated with the first call to Info, Error
// or Enabled or the first call to those function on any child loggers created via V, WithName or
// WithValues.
func NewTextLog() logger.Logger {
	return &sink{core: &core{
		w:        os.Stderr,
		colorize: true,
		//humanize: true,
		tsFormat: "15:04:05.000000",
	}}
}

// sink is a logger sink that writes messages in the logfmt style.
// See https://www.brandur.org/logfmt for more information.
type sink struct {
	core *core
}

// Info logs a non-error message with the given key/value pairs as context.
func (l *sink) Info(level uint, msg string, kvs ...logger.Tagger) {
	var ff = logfmt.NewEncoder(l.core.w)
	ff.EncodeKeyval("level", level)
	ff.EncodeKeyval("prefix", l.core.name)
	ff.EncodeKeyval("msg", msg)
	for i := range kvs {
		ff.EncodeKeyval(kvs[i].Key(), kvs[i].Value())
	}
	return
	l.core.write(level, "info", msg, l.core.flatten(kvs...))
}

// Error logs an error, with the given message and key/value pairs as context.
func (l *sink) Error(err error, msg string, kvs ...logger.Tagger) {
	kvs = append(kvs, logger.Tag("error", err))
	l.core.write(0, "error", msg, l.core.flatten())
}

type core struct {
	w        io.Writer
	name     string
	tsFormat string

	humanize bool
	colorize bool
}

func (c *core) write(level uint, humanprefix, msg string, values string, extras ...logger.Tagger) {
	var b bytes.Buffer
	if c.humanize {
		if c.colorize {
			if humanprefix == "error" {
				humanprefix = colorRed + humanprefix + colorDefault
			} else {
				humanprefix = colorGreen + humanprefix + " " + colorDefault
			}
		}

		b.WriteString(fmt.Sprintf("%d %-5s | %15s | %-30s", level, humanprefix, time.Now().UTC().Format("15:04:05.000000"), msg))
		if c.name != "" {
			b.WriteRune(' ')
			b.WriteString(c.key("logger"))
			b.WriteString("=")
			b.WriteString(c.name)
		}
	} else {
		b.WriteString("level=")
		b.WriteString(strconv.Itoa(int(level)))
		if c.name != "" {
			b.WriteRune(' ')
			b.WriteString("logger=")
			b.WriteString(quote(c.name))
		}
		if c.tsFormat != "" {
			b.WriteRune(' ')
			b.WriteString("ts=")
			b.WriteString(quote(time.Now().UTC().Format(c.tsFormat)))
		}
		b.WriteRune(' ')
		b.WriteString("msg=")
		b.WriteString(quote(msg))
	}

	if len(extras) > 0 {
		b.WriteRune(' ')
		b.WriteString(c.flatten(extras...))
	}

	if values != "" {
		b.WriteRune(' ')
		b.WriteString(values)
	}
	b.WriteRune('\n')
	_, _ = c.w.Write(b.Bytes())
}

func (c *core) flatten(kvs ...logger.Tagger) string {
	if len(kvs) == 0 {
		return ""
	}

	var b strings.Builder
	for i := 0; i < len(kvs); i++ {
		if i > 0 {
			b.WriteRune(' ')
		}

		b.WriteString(c.key(stringify(kvs[i].Key())))
		b.WriteRune('=')
		b.WriteString(stringify(kvs[i].Value()))
	}

	return b.String()
}

func (c *core) key(s string) string {
	if !c.colorize {
		return s
	}

	switch s {
	case "error":
		return colorRed + s + colorDefault
	case "logger", "caller":
		return colorBlue + s + colorDefault
	default:
		return colorYellow + s + colorDefault
	}
}

func stringify(v interface{}) string {
	var s string
	switch vv := v.(type) {
	case string:
		s = vv
	case fmt.Stringer:
		s = vv.String()
	case error:
		s = vv.Error()
	default:
		s = fmt.Sprint(v)
	}
	return quote(s)
}

func quote(s string) string {
	if strings.ContainsAny(s, " ") {
		return fmt.Sprintf("%q", s)
	}
	return s
}

const (
	colorDefault = "\x1b[0m"
	colorRed     = "\x1b[1;31m"
	colorGreen   = "\x1b[1;32m"
	colorYellow  = "\x1b[1;33m"
	colorBlue    = "\x1b[1;34m"
)
