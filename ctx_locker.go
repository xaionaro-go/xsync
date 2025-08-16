package xsync

import (
	"context"

	"github.com/facebookincubator/go-belt/tool/logger"
)

// TODO: move to a separate package
type CtxLocker chan struct{}

func (l CtxLocker) ManualLock(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		loggingEnabled := IsLoggingEnabled(ctx)
		if loggingEnabled {
			logger.Tracef(ctx, "context is closed")
		}
		return false
	case l <- struct{}{}:
		return true
	}
}

func (l CtxLocker) ManualTryLock(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		loggingEnabled := IsLoggingEnabled(ctx)
		if loggingEnabled {
			logger.Tracef(ctx, "context is closed")
		}
		return false
	case l <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l CtxLocker) ManualUnlock() {
	select {
	case <-l:
	default:
		panic("not locked!")
	}
}

func (l CtxLocker) Do(
	ctx context.Context,
	fn func(),
) {
	l.ManualLock(ctx)
	defer l.ManualUnlock()
	fn()
}
