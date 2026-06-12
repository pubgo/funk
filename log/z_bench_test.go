package log_test

import (
	"context"
	"io"
	"testing"

	"github.com/pubgo/funk/v2/log"
)

func BenchmarkInfoWithContext(b *testing.B) {
	logger := log.Output(io.Discard).WithFields(log.Fields{"component": "bench", "scope": "logger"})
	ctx := log.CreateFieldsCtx(context.Background(), log.Fields{"scope": "ctx", "request_id": "req-bench"})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx).Str("iteration", "bench").Msg("benchmark")
	}
}

func BenchmarkUpdateFieldsCtx(b *testing.B) {
	ctx := log.CreateFieldsCtx(context.Background(), log.Fields{"request_id": "req-bench", "trace_id": "trace-bench"})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = log.UpdateFieldsCtx(ctx, log.Fields{"user_id": i, "attempt": i % 3})
	}
}
