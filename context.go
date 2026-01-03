// context.go provides helper functions to manage xsync settings (like logging and deadlock detection) via context.

package xsync

import (
	"context"
)

type CtxKeyLogging struct{}

// WithNoLogging sets whether logging is disabled.
// DEPRECATED: use WithLogging
func WithNoLogging(ctx context.Context, noLogging bool) context.Context {
	return context.WithValue(ctx, CtxKeyLogging{}, !noLogging)
}

// IsNoLogging returns whether logging is disabled.
// DEPRECATED: use IsLoggingEnabled
func IsNoLogging(ctx context.Context) bool {
	v, _ := ctx.Value(CtxKeyLogging{}).(bool)
	return !v
}

func WithLoggingEnabled(ctx context.Context, logging bool) context.Context {
	return context.WithValue(ctx, CtxKeyLogging{}, logging)
}

func IsLoggingEnabled(ctx context.Context) bool {
	v, _ := ctx.Value(CtxKeyLogging{}).(bool)
	return v
}

type CtxKeyEnableDeadlock struct{}

// WithEnableDeadlock sets whether the deadlock detector is enabled.
//
// DEPRECATED: use WithDeadlockDetectorEnabled
func WithEnableDeadlock(ctx context.Context, enableDeadlock bool) context.Context {
	return WithDeadlockDetectorEnabled(ctx, enableDeadlock)
}

// IsEnableDeadlock returns whether the deadlock detector is enabled.
//
// DEPRECATED: use IsDeadlockDetectorEnabled
func IsEnableDeadlock(ctx context.Context) bool {
	return IsDeadlockDetectorEnabled(ctx)
}

func WithDeadlockDetectorEnabled(ctx context.Context, enableDeadlock bool) context.Context {
	return context.WithValue(ctx, CtxKeyEnableDeadlock{}, enableDeadlock)
}

func IsDeadlockDetectorEnabled(ctx context.Context) bool {
	v, ok := ctx.Value(CtxKeyEnableDeadlock{}).(bool)
	if !ok {
		return false
	}
	return v
}

type CtxKeyAllowUnlockNotLocked struct{}

func WithAllowUnlockNotLocked(ctx context.Context, allow bool) context.Context {
	return context.WithValue(ctx, CtxKeyAllowUnlockNotLocked{}, allow)
}

func IsAllowUnlockNotLocked(ctx context.Context) bool {
	v, ok := ctx.Value(CtxKeyAllowUnlockNotLocked{}).(bool)
	if !ok {
		return false
	}
	return v
}
