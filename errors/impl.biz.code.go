package errors

func (t *baseErr) BizCode() string {
	return t.bizCode
}

func (t *baseErr) AddBizCode(biz string) {
	t.bizCode = biz
}
