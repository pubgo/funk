package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	"github.com/pubgo/funk/v2/env"
	"github.com/pubgo/funk/v2/errors"
	"github.com/pubgo/funk/v2/log"
)

// celEngine provides CEL expression evaluation with security features
type celEngine struct {
	celEnv     *cel.Env
	workDir    string
	envSpecMap EnvSpecMap // allowed env vars from patch_envs
}

// newCelEngine creates a new CEL engine with built-in functions
func newCelEngine(cfg *config) (*celEngine, error) {
	// Define custom functions
	embedFunc := cel.Function("embed",
		cel.Overload("embed_string",
			[]*cel.Type{cel.StringType},
			cel.StringType,
			cel.UnaryBinding(func(arg ref.Val) ref.Val {
				name, ok := arg.Value().(string)
				if !ok {
					return types.NewErr("embed: expected string argument")
				}
				if name == "" {
					return types.String("")
				}

				// Security: validate path to prevent path traversal attacks
				safePath, err := securePath(cfg.workDir, name)
				if err != nil {
					log.Error().Err(err).
						Str("name", name).
						Str("workDir", cfg.workDir).
						Msg("embed: path security check failed")
					return types.String("")
				}

				d, err := os.ReadFile(safePath)
				if err != nil {
					log.Error().Err(err).
						Str("path", safePath).
						Msg("embed: failed to read file")
					return types.String("")
				}

				return types.String(strings.TrimSpace(base64.StdEncoding.EncodeToString(d)))
			}),
		),
	)

	configDirFunc := cel.Function("config_dir",
		cel.Overload("config_dir_void",
			[]*cel.Type{},
			cel.StringType,
			cel.FunctionBinding(func(args ...ref.Val) ref.Val {
				return types.String(cfg.workDir)
			}),
		),
	)

	// env function - get env var, must be defined in patch_envs
	envFunc := cel.Function("env",
		cel.Overload("env_string",
			[]*cel.Type{cel.StringType},
			cel.StringType,
			cel.UnaryBinding(func(arg ref.Val) ref.Val {
				key, ok := arg.Value().(string)
				if !ok {
					return types.NewErr("env: expected string argument")
				}

				// Check if env var is defined in envSpecMap
				if cfg.envSpecMap != nil {
					key = strings.TrimSpace(strings.ToUpper(key))
					if _, defined := cfg.envSpecMap[key]; !defined {
						return types.NewErr("env: variable %q is not defined in patch_envs, all env vars must be declared", key)
					}
				}
				return types.String(env.Get(key))
			}),
		),
	)

	// envs function - return all defined env vars as map
	envsFunc := cel.Function("envs",
		cel.Overload("envs_void",
			[]*cel.Type{},
			cel.MapType(cel.StringType, cel.StringType),
			cel.FunctionBinding(func(args ...ref.Val) ref.Val {
				result := make(map[string]string)
				if cfg.envSpecMap != nil {
					for name, spec := range cfg.envSpecMap {
						result[name] = spec.GetValue()
					}
				}
				return types.DefaultTypeAdapter.NativeToValue(result)
			}),
		),
	)

	// Build CEL environment options (functions only, no variables)
	opts := []cel.EnvOption{
		embedFunc,
		configDirFunc,
		envFunc,
		envsFunc,
	}

	// Add custom registered functions
	for name, fn := range globalManager.GetExprFuncs() {
		fnOpt, err := createCelFunction(name, fn)
		if err != nil {
			log.Warn().Err(err).Str("name", name).Msg("failed to register custom CEL function")
			continue
		}
		opts = append(opts, fnOpt)
	}

	celEnv, err := cel.NewEnv(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CEL environment")
	}

	return &celEngine{
		celEnv:     celEnv,
		workDir:    cfg.workDir,
		envSpecMap: cfg.envSpecMap,
	}, nil
}

// Eval evaluates a CEL expression and returns the result
func (e *celEngine) Eval(expression string) (any, error) {
	expression = strings.TrimSpace(expression)

	// Parse and check the expression
	ast, issues := e.celEnv.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return nil, errors.Wrapf(issues.Err(), "CEL compile error for expression: %q", expression)
	}

	// Create the program
	prg, err := e.celEnv.Program(ast)
	if err != nil {
		return nil, errors.Wrapf(err, "CEL program error for expression: %q", expression)
	}

	// Evaluate (no input variables, only functions)
	out, _, err := prg.Eval(map[string]any{})
	if err != nil {
		return nil, errors.Wrapf(err, "CEL eval error for expression: %q", expression)
	}

	return out.Value(), nil
}

// createCelFunction creates a CEL function option from a Go function
// Supported signatures:
//   - fn() T
//   - fn() (T, error)
//   - fn() error
//   - fn(T) R
//   - fn(T) (R, error)
//   - fn(T) error
func createCelFunction(name string, fn any) (cel.EnvOption, error) {
	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()

	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("expected function, got %T", fn)
	}

	numIn := fnType.NumIn()
	numOut := fnType.NumOut()

	if numOut < 1 || numOut > 2 {
		return nil, fmt.Errorf("function must have 1 or 2 return values, got %d", numOut)
	}
	if numIn > 1 {
		return nil, fmt.Errorf("function must have 0 or 1 input parameter, got %d", numIn)
	}

	// Check if last return value is error
	hasError := numOut >= 1 && fnType.Out(numOut-1).Implements(reflect.TypeOf((*error)(nil)).Elem())

	// Determine the result type (non-error return type)
	var resultType *cel.Type
	if hasError && numOut == 1 {
		// fn() error or fn(T) error - returns nothing useful, use NullType/DynType
		resultType = cel.DynType
	} else if hasError && numOut == 2 {
		// fn() (T, error) or fn(T) (R, error)
		resultType = goCelType(fnType.Out(0))
	} else {
		// fn() T or fn(T) R
		resultType = goCelType(fnType.Out(0))
	}

	// Create appropriate overload based on function signature
	switch numIn {
	case 0:
		// func() T, func() (T, error), or func() error
		return cel.Function(name,
			cel.Overload(name+"_void",
				[]*cel.Type{},
				resultType,
				cel.FunctionBinding(func(args ...ref.Val) ref.Val {
					results := fnVal.Call(nil)
					return handleFuncResults(results, hasError, numOut)
				}),
			),
		), nil

	case 1:
		// func(T) R, func(T) (R, error), or func(T) error
		return cel.Function(name,
			cel.Overload(name+"_unary",
				[]*cel.Type{goCelType(fnType.In(0))},
				resultType,
				cel.UnaryBinding(func(arg ref.Val) ref.Val {
					goArg := celToGoValue(arg, fnType.In(0))
					results := fnVal.Call([]reflect.Value{reflect.ValueOf(goArg)})
					return handleFuncResults(results, hasError, numOut)
				}),
			),
		), nil

	default:
		return nil, fmt.Errorf("unsupported function signature: %v", fnType)
	}
}

// handleFuncResults processes function return values and converts to CEL value
func handleFuncResults(results []reflect.Value, hasError bool, numOut int) ref.Val {
	if hasError {
		// Check error (always last return value)
		errVal := results[numOut-1]
		if !errVal.IsNil() {
			err := errVal.Interface().(error)
			return types.NewErr("%s", err.Error())
		}

		// fn() error or fn(T) error - return null on success
		if numOut == 1 {
			return types.NullValue
		}

		// fn() (T, error) or fn(T) (R, error) - return first value
		return goToCelValue(results[0].Interface())
	}

	// fn() T or fn(T) R - return the single value
	return goToCelValue(results[0].Interface())
}

// goCelType converts Go type to CEL type
func goCelType(t reflect.Type) *cel.Type {
	switch t.Kind() {
	case reflect.String:
		return cel.StringType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cel.IntType
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return cel.UintType
	case reflect.Float32, reflect.Float64:
		return cel.DoubleType
	case reflect.Bool:
		return cel.BoolType
	default:
		return cel.DynType
	}
}

// goToCelValue converts Go value to CEL value
func goToCelValue(v any) ref.Val {
	if v == nil {
		return types.NullValue
	}
	switch val := v.(type) {
	case string:
		return types.String(val)
	case int:
		return types.Int(val)
	case int64:
		return types.Int(val)
	case float64:
		return types.Double(val)
	case bool:
		return types.Bool(val)
	default:
		return types.DefaultTypeAdapter.NativeToValue(v)
	}
}

// celToGoValue converts CEL value to Go value
func celToGoValue(v ref.Val, targetType reflect.Type) any {
	goVal := v.Value()

	// Type conversion if needed
	switch targetType.Kind() {
	case reflect.String:
		if s, ok := goVal.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", goVal)
	case reflect.Int, reflect.Int64:
		if i, ok := goVal.(int64); ok {
			return i
		}
	case reflect.Float64:
		if f, ok := goVal.(float64); ok {
			return f
		}
	case reflect.Bool:
		if b, ok := goVal.(bool); ok {
			return b
		}
	}

	return goVal
}
