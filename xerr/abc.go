package xerr

type XErr interface {
	xErr()
	Error() string
	String() string
	DebugPrint()
	Unwrap() error
	Wrap(args ...interface{}) XErr
	WrapKV(k string, v interface{}) XErr
	WrapF(msg string, args ...interface{}) XErr
}
