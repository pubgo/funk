package logger

type Map map[string]any

type Valuer func() (BytesL, error)
type Fields []Field
type Field interface {
	Name() string
	Kind() FieldType
	Value() Valuer
}

type FieldType int

const (
	_ FieldType = iota
	BytesType
	RawType
)
