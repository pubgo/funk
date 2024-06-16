package template

import (
	"bytes"
	"io/fs"
	"sync"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/log"
)

type Template struct {
	fs        fs.FS
	templates sync.Map
	globals   native.Declarations
}

func Build[Data any](t *Template, name string, data *Data) (string, error) {
	_, ok := t.templates.Load(name)
	if !ok {
		vars := make(native.Declarations)
		for k := range t.globals {
			vars[k] = t.globals[k]
		}

		vars["data"] = (*Data)(nil)
		template, err := scriggo.BuildTemplate(t.fs, name, &scriggo.BuildOptions{Globals: vars})
		if err != nil {
			log.Err(err).Str("name", name).Any("data", data).Msg("failed to build template")
			return "", err
		}
		t.templates.Store(name, template)
	}

	tt, _ := t.templates.Load(name)
	var buf bytes.Buffer
	err := tt.(*scriggo.Template).Run(&buf, map[string]interface{}{"data": data}, nil)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func New(fs fs.FS, globals native.Declarations) *Template {
	assert.If(fs == nil, "fs should not be nil")

	if globals == nil {
		globals = make(native.Declarations)
	}

	return &Template{fs: fs, globals: globals}
}
