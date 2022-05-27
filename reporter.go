package rdy

import (
	"context"
)

type Reporter interface {
	L1(ctx context.Context, format string, args ...interface{})
	L2(ctx context.Context, format string, args ...interface{})
}

type safeReporter struct {
	Reporter
}

func (s *safeReporter) L1(ctx context.Context, format string, args ...interface{}) {
	if s.Reporter != nil {
		s.Reporter.L1(ctx, format, args...)
	}
}

func (s *safeReporter) L2(ctx context.Context, format string, args ...interface{}) {
	if s.Reporter != nil {
		s.Reporter.L2(ctx, format, args...)
	}
}
