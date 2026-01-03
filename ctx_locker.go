// ctx_locker.go implements CtxLocker, a channel-based synchronization primitive that supports context-aware locking.

package xsync

import (
	"context"

	"github.com/facebookincubator/go-belt/tool/logger"
)

// CtxLocker is a channel-based synchronization primitive that supports context-aware locking.
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

func (l CtxLocker) ManualUnlock(ctx context.Context) {
	select {
	case <-l:
	default:
		if IsAllowUnlockNotLocked(ctx) {
			return
		}
		panic("not locked!")
	}
}

func (l CtxLocker) Do(
	ctx context.Context,
	fn func(),
) {
	if !l.ManualLock(ctx) {
		return
	}
	defer l.ManualUnlock(ctx)
	fn()
}
