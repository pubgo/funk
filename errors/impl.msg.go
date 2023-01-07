package errors

func (t *baseErr) Msg() string {
	return t.msg
}

func (t *baseErr) AddMsg(msg string) {
	t.msg = msg
}
