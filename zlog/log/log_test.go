package log

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	t.Run("no level", func(t *testing.T) {
		var log = New(Config{
			Level: "",
		})
		assert.Equal(t, zerolog.InfoLevel, log.GetLevel())
	})

	t.Run("debug level", func(t *testing.T) {
		var log = New(Config{
			Level: "debug",
		})
		assert.Equal(t, zerolog.DebugLevel, log.GetLevel())
	})

	t.Run("module field", func(t *testing.T) {
		b := strings.Builder{}
		l := New(Config{
			Level:  "debug",
			AsJson: true,
			Writer: &b,
		})

		const message = "test"
		l = Module(&l, message)
		l.Debug().Msg(time.Now().String())

		var log map[string]interface{}
		assert.NoError(t, json.Unmarshal([]byte(b.String()), &log))
		assert.Equal(t, message, log["module"])
	})

	t.Run("json format", func(t *testing.T) {
		b := strings.Builder{}
		l := New(Config{
			Level:  "debug",
			AsJson: true,
			Writer: &b,
		})

		const message = "Should appear as json"
		l.Debug().Msg(message)

		var log map[string]interface{}
		assert.NoError(t, json.Unmarshal([]byte(b.String()), &log))
		assert.Equal(t, message, log["message"])
	})

	t.Run("text format", func(t *testing.T) {
		b := strings.Builder{}
		l := New(Config{
			Level:  "info",
			AsJson: false,
			Writer: &b,
		})

		const (
			message1 = "Should not appear"
			message2 = "Should appear as text"
		)

		l.Debug().Msg(message1)
		l.Info().Msg(message2)

		lines := strings.Split(strings.Trim(b.String(), "\n"), "\n")
		assert.Equal(t, 1, len(lines))

		var log map[string]interface{}
		assert.Error(t, json.Unmarshal([]byte(lines[0]), &log))

		assert.True(t, strings.Contains(lines[0], message2))
	})
}
