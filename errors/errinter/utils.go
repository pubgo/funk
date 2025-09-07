package errinter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/k0kubun/pp/v3"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/pubgo/funk"
	"github.com/pubgo/funk/log/logutil"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/proto/errorpb"
)

func ParseError(val interface{}) error {
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
		return errors.New(SimplePrint(v))
	}
}

var Simple = sync.OnceValue(func() *pp.PrettyPrinter {
	printer := pp.New()
	printer.SetColoringEnabled(false)
	printer.SetExportedOnly(false)
	printer.SetOmitEmpty(true)
	printer.SetMaxDepth(3)
	return printer
})

func SimplePrint(v interface{}) string {
	return strings.ReplaceAll(Simple().Sprint(v), "\n", "")
}

func GetErrorId(err error) string {
	if err == nil {
		return ""
	}

	for err != nil {
		if v, ok := err.(Error); ok {
			return v.ID()
		}

		err = Unwrap(err)
	}

	return ""
}

func Unwrap(err error) error {
	u, ok := err.(ErrUnwrap)
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func ParseErrToPb(err error) proto.Message {
	switch err1 := err.(type) {
	case nil:
		return nil
	case ErrorProto:
		return err1.Proto()
	case GRPCStatus:
		return err1.GRPCStatus().Proto()
	case proto.Message:
		return err1
	default:
		return &errorpb.ErrMsg{Msg: err.Error(), Detail: fmt.Sprintf("%v", err)}
	}
}

func Debug(err error) {
	if err == nil {
		return
	}

	if _err, ok := err.(fmt.Stringer); ok {
		_, _ = fmt.Fprintln(os.Stderr, _err.String())
		return
	}

	pretty.SetDefaultMaxDepth(20)
	pretty.Println(err)
}

func MustTagsToAny(tags ...*errorpb.Tag) []*anypb.Any {
	if len(tags) == 0 {
		return nil
	}

	return lo.Map(tags, func(item *errorpb.Tag, index int) *anypb.Any {
		return MustProtoToAny(item)
	})
}

func MustStructToAny(p map[string]any) *anypb.Any {
	if p == nil {
		return nil
	}

	pb, err := structpb.NewStruct(p)
	if err != nil {
		log.Err(err).
			Any("params", p).
			Func(logutil.WithNotice()).
			Stack().
			Msgf("failed to encode map-any to struct protobuf")
		return anyToProtobuf(p)
	} else {
		return MustProtoToAny(pb)
	}
}

func anyToProtobuf(v any) *anypb.Any {
	data, err := json.Marshal(v)
	if err != nil {
		return &anypb.Any{
			TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
			Value:   []byte(fmt.Sprintf("err:%s detail:%#v", err.Error(), v)),
		}
	}
	return &anypb.Any{
		TypeUrl: "type.googleapis.com/google.protobuf.Struct",
		Value:   data,
	}
}

func MustProtoToAny(p proto.Message) *anypb.Any {
	switch p := p.(type) {
	case nil:
		return nil
	case *anypb.Any:
		return p
	}

	pb, err := anypb.New(p)
	if err != nil {
		params := prototext.Format(p)
		log.Err(err).
			Str("protobuf", params).
			Func(logutil.WithNotice()).
			Stack().
			Msgf("failed to encode protobuf message to any protobuf")
		pb, _ = anypb.New(structpb.NewStringValue(params))
		return pb
	}

	return pb
}
