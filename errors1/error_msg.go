package error1

import "github.com/pubgo/funk/proto/errorpb"

func New(msg string) *ErrMsg {
	return &ErrMsg{
		msg: &errorpb.ErrMsg{Msg: msg, Tags: make(map[string]string)},
	}
}

type ErrMsg struct {
	msg *errorpb.ErrMsg
}

func (e *ErrMsg) clone() *ErrMsg {
	return &ErrMsg{
		msg: e.msg.DeepCopy(),
	}
}

func (e *ErrMsg) WithMsg(msg string) *ErrMsg {
	var c = e.clone()
	c.msg.Msg = msg
	return c
}

func (e *ErrMsg) WithTags(tags map[string]string) *ErrMsg {
	var c = e.clone()
	if tags == nil || len(tags) == 0 {
		return c
	}

	for k, v := range tags {
		c.msg.Tags[k] = v
	}

	return c
}

func (e *ErrMsg) WithTag(key, value string) *ErrMsg {
	var c = e.clone()
	if c.msg.Tags == nil {
		c.msg.Tags = make(map[string]string)
	}
	c.msg.Tags[key] = value
	return c
}

func (e *ErrMsg) Error() string {
	return e.msg.Msg
}

func (e *ErrMsg) String() string {
	return e.msg.String()
}

func (e *ErrMsg) MarshalJSON() ([]byte, error) {
	return e.msg.MarshalJSON()
}

func (e *ErrMsg) Proto() *errorpb.ErrMsg {
	return e.msg
}
