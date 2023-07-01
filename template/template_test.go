package template

import (
	"testing"

	"github.com/open2b/scriggo"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	fs := scriggo.Files{"index.txt": []byte(`Hello {{ data.Hello }}`)}
	tt := New(fs, nil)

	type Data struct {
		Hello string
	}

	data, err := Build(tt, "index.txt", &Data{Hello: "World"})
	assert.Nil(t, err)
	t.Log(data)
	assert.Equal(t, data, `Hello World`)
}
