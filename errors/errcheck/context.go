package errcheck

import (
	"context"
)

type checkCtx struct{}

type ErrChecker func(context.Context, error) error

func CreateCtx(ctx context.Context, errChecks []ErrChecker) context.Context {
	if ctx == nil {
		ctx = context.Background()
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
