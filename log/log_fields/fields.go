package log_fields

import (
	"encoding/base64"
	"strconv"
	"time"

	jjson "github.com/goccy/go-json"
	"github.com/pubgo/funk/convert"
	"github.com/pubgo/funk/log/log_config"
	"github.com/pubgo/funk/log/logutil"
	"github.com/pubgo/funk/logger"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type field struct {
	key       string
	fieldType logger.FieldType
	value     logger.Valuer
}

func (f field) Name() string {
	return f.key
}

func (f field) Kind() logger.FieldType {
	return f.fieldType
}

func (f field) Value() logger.Valuer {
	return f.value
}

func Timestamp(t time.Time) logger.Valuer {
	return func() (logger.BytesL, error) {
		return logger.BytesL{logutil.AppendTime(nil, t, time.RFC3339)}, nil
	}
}

func Hex(key string, val []byte) logger.Field {
	return field{key: key, fieldType: logger.BytesType, value: func() (logger.BytesL, error) {
		return logger.BytesL{logutil.AppendHex(nil, val)}, nil
	}}
}

func Base64(key string, val []byte) logger.Field {
	return field{key: key, fieldType: logger.BytesType, value: func() (logger.BytesL, error) {
		buf := make([]byte, base64.StdEncoding.EncodedLen(len(val)))
		base64.StdEncoding.Encode(buf, val)
		return logger.BytesL{buf}, nil
	}}
}

func Bool(key string, val bool) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		return logger.BytesL{strconv.AppendBool(nil, val)}, nil
	}}
}

func BoolL(key string, val ...bool) logger.Field {
	return field{key: key, fieldType: logger.RawArrayType, value: func() (logger.BytesL, error) {
		var data = make(logger.BytesL, len(val))
		for i := range val {
			data[i] = strconv.AppendBool(nil, val[i])
		}
		return data, nil
	}}
}

func Float64(key string, val float64) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		return logger.BytesL{logutil.AppendFloat64(nil, val)}, nil
	}}
}

func Float64L(key string, val ...float64) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		var data = make(logger.BytesL, len(val))
		for i := range val {
			data[i] = logutil.AppendFloat64(nil, val[i])
		}
		return data, nil
	}}
}

func Int(key string, val int) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		return logger.BytesL{logutil.AppendInt(nil, val)}, nil
	}}
}

func IntL(key string, val ...int) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		var data = make(logger.BytesL, len(val))
		for i := range val {
			data[i] = logutil.AppendInt(nil, val[i])
		}
		return data, nil
	}}
}

func Int64(key string, val int64) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		return logger.BytesL{logutil.AppendInt64(nil, val)}, nil
	}}
}

func Int64L(key string, val ...int64) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		var data = make(logger.BytesL, len(val))
		for i := range val {
			data[i] = logutil.AppendInt64(nil, val[i])
		}
		return data, nil
	}}
}

func Uint(key string, val uint) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		return logger.BytesL{logutil.AppendUint(nil, val)}, nil
	}}
}

func UintL(key string, val ...uint) logger.Field {
	return field{key: key, fieldType: logger.RawArrayType, value: func() (logger.BytesL, error) {
		var data = make(logger.BytesL, len(val))
		for i := range val {
			data[i] = logutil.AppendUint(nil, val[i])
		}
		return data, nil
	}}
}

func Uint64(key string, val uint64) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		return logger.BytesL{logutil.AppendUint64(nil, val)}, nil
	}}
}

func Uint64L(key string, val ...uint64) logger.Field {
	return field{key: key, fieldType: logger.RawArrayType, value: func() (logger.BytesL, error) {
		var data = make(logger.BytesL, len(val))
		for i := range val {
			data[i] = logutil.AppendUint64(nil, val[i])
		}
		return data, nil
	}}
}

func String(key string, val string) logger.Field {
	return field{key: key, fieldType: logger.BytesType, value: func() (logger.BytesL, error) {
		return logger.BytesL{convert.S2B(val)}, nil
	}}
}

func StringL(key string, val ...string) logger.Field {
	return field{key: key, fieldType: logger.BytesArrayType, value: func() (logger.BytesL, error) {
		var data = make(logger.BytesL, len(val))
		for i := range val {
			data[i] = convert.S2B(val[i])
		}
		return data, nil
	}}
}

func JsonRaw(key string, val []byte) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		return logger.BytesL{val}, nil
	}}
}

func Json(key string, val interface{}) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		if val == nil {
			return logger.BytesL{[]byte("null")}, nil
		}

		var data, err = jjson.Marshal(val)
		if err != nil {
			return nil, err
		}
		return logger.BytesL{data}, nil
	}}
}

func JsonL(key string, val ...interface{}) logger.Field {
	return field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		if val == nil || len(val) == 0 {
			return logger.BytesL{[]byte("null")}, nil
		}

		var data, err = jjson.Marshal(val)
		if err != nil {
			return nil, err
		}
		return logger.BytesL{data}, nil
	}}
}

func Error(err error) []logger.Field {
	if err == nil {
		return nil
	}

	return []logger.Field{
		field{key: log_config.FieldErrorKey, fieldType: logger.BytesType,
			value: func() (logger.BytesL, error) {
				return logger.BytesL{}, nil
			}},
		field{key: log_config.FieldErrorDetailKey, fieldType: logger.BytesType,
			value: func() (logger.BytesL, error) {
				return logger.BytesL{}, nil
			}},
	}
}

func Time(key string, val time.Time) logger.Field {
	if len(log_config.TimeFormat) > 0 {
		return String(key, val.UTC().Format(log_config.TimeFormat))
	}

	return Int64(key, val.UnixMilli())
}

func Dur(key string, val time.Duration) logger.Field {
	return Int64(key, val.Milliseconds())
}

func Duration(val time.Duration) logger.Field {
	return Int64("duration", val.Milliseconds())
}

var jsonMarshaller = &protojson.MarshalOptions{EmitUnpopulated: true}

func Protobuf(key string, val proto.Message) logger.Field {
	return &field{key: key, fieldType: logger.RawType, value: func() (logger.BytesL, error) {
		var data, err = jsonMarshaller.Marshal(val)
		if err != nil {
			return nil, err
		}
		return logger.BytesL{data}, nil
	}}
}
