package errparser

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/pubgo/funk/v2"
	"github.com/pubgo/funk/v2/stack"
)

func Parse(val any) error { return parseError(val) }

type ErrParser = func(val any) (error, bool)

var errParserRegistry = make(map[string]ErrParser)

func RegisterParser(errParser ErrParser) bool {
	key := stack.CallerWithFunc(errParser).String()
	if errParserRegistry[key] != nil {
		slog.Error("errParser already exists", "parser", errParserRegistry[key])
		return false
	}

	if errParser == nil {
		slog.Error("errParser is nil")
		return false
	}

	errParserRegistry[key] = errParser
	return true
}

func parseError(val any) error {
	if funk.IsNil(val) {
		return nil
	}

	switch v := val.(type) {
	case nil:
		return nil
	case error:
		return v
	case string:
		return errors.New(v)
	case []byte:
		return errors.New(string(v))
	default:
		for _, value := range errParserRegistry {
			parsedErr, ok := value(v)
			if ok && parsedErr != nil {
				return parsedErr
			}
		}

		return fmt.Errorf("%#v", val)
	}
}
