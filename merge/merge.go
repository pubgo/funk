package merge

import (
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/pretty"
	"github.com/pubgo/funk/result"
)

type Option func(opts *copier.Option)

// Copy
// struct<->struct
// 各种类型结构体之间的field copy
func Copy[A any, B any](dst A, src B, opts ...Option) result.Result[A] {
	var opt = copier.Option{DeepCopy: true, IgnoreEmpty: true}
	for i := range opts {
		opts[i](&opt)
	}

	var errH = func(err error) error {
		return errors.WrapEventFn(err, func(evt *errors.Event) {
			evt.Any("dst", dst)
			evt.Any("src", src)
			evt.Any("decoder_config", pretty.Sprint(opt))
		})
	}

	var err = copier.CopyWithOption(dst, src, opt)
	if err != nil {
		return result.Err[A](errH(err))
	}

	return result.OK(dst)
}

func Struct[A any, B any](dst A, src B, opts ...Option) result.Result[A] {
	return Copy(dst, src, opts...)
}

// MapStruct
// map<->struct
// map和结构体相互转化
func MapStruct[A any, B any](dst A, src B, opts ...func(cfg *mapstructure.DecoderConfig)) (r result.Result[A]) {
	var cfg = &mapstructure.DecoderConfig{
		TagName:          "json",
		Metadata:         nil,
		Result:           &dst,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}

	for i := range opts {
		opts[i](cfg)
	}

	var errH = func(err error) error {
		return errors.WrapEventFn(err, func(evt *errors.Event) {
			evt.Any("dst", dst)
			evt.Any("src", src)
			evt.Any("decoder_config", pretty.Sprint(cfg))
		})
	}

	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return r.WithErr(errH(err))
	}

	err = decoder.Decode(src)
	if err != nil {
		return r.WithErr(errH(err))
	}

	return result.OK(dst)
}
