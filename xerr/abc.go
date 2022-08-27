package xerr

type XErr interface {
	xErr()
	Error() string
	String() string
	DebugPrint()
	Stack() string
	Unwrap() error
	Cause() error
	Wrap(args ...interface{}) XErr
	WithMeta(k string, v interface{}) XErr
	WrapF(msg string, args ...interface{}) XErr
}
