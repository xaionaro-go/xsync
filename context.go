package xsync

import (
	"context"
)

type CtxKeyLogging struct{}

// DEPRECATED: use WithLogging
func WithNoLogging(ctx context.Context, noLogging bool) context.Context {
	return context.WithValue(ctx, CtxKeyLogging{}, !noLogging)
}

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

// DEPRECATED: use WithDeadlockDetectorEnabled
func WithEnableDeadlock(ctx context.Context, enableDeadlock bool) context.Context {
	return WithDeadlockDetectorEnabled(ctx, enableDeadlock)
}

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
