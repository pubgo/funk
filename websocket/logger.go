package websocket

var logger = wlog.DefaultLogger()

func SetLogger(l *wlog.Logger) {
	logger = l.WithFields(map[string]interface{}{"pkg": "websocket", "type": "internal"})
}

type LWriter struct{}

func (writer LWriter) Write(p []byte) (n int, err error) {
	logger.Errorf(string(p))
	return len(p), nil
}
