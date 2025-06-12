package errcheck

import (
	"context"

	"github.com/samber/lo"
)

type checkCtx struct{}

type ErrChecker func(context.Context, error) error

func CreateCtx(ctx context.Context, errChecks []ErrChecker, upsert ...bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if lo.FirstOrEmpty(upsert) {
		checkers, ok := ctx.Value(checkCtx{}).([]ErrChecker)
		if ok {
			errChecks = append(errChecks, checkers...)
		}
	}

	return context.WithValue(ctx, checkCtx{}, errChecks)
}

func GetCheckersFromCtx(ctx context.Context) []ErrChecker {
	if ctx == nil {
		return nil
	}

	checkers, ok := ctx.Value(checkCtx{}).([]ErrChecker)
	if !ok {
		return nil
	}

	return checkers
}
