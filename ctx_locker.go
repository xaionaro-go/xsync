package xsync

import (
	"context"

	"github.com/facebookincubator/go-belt/tool/logger"
)

// TODO: move to a separate package
type CtxLocker chan struct{}

func (l CtxLocker) Lock(ctx context.Context) bool {
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

func (l CtxLocker) TryLock(ctx context.Context) bool {
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

func (l CtxLocker) Unlock() {
	select {
	case <-l:
	default:
		panic("not locked!")
	}
}
