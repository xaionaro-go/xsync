// context.go provides helper functions to manage xsync settings (like logging and deadlock detection) via context.

package xsync

import (
	"context"
	"sync/atomic"
)

// loggingPossiblyEnabled is the master switch for IsLoggingEnabled. While
// false (the default), IsLoggingEnabled short-circuits to false without
// performing a context.Value chain walk on the hot RWMutex Lock/Unlock path.
// The gate is auto-flipped to true the first time a caller invokes one of
// the per-context setters (WithLoggingEnabled / WithLogging / WithNoLogging)
// with a value that asks for logging on, so callers do not need to remember
// to call SetLoggingPossiblyEnabled themselves. The flip is monotonic-on:
// asking-false never flips the gate down, because other live contexts may
// still want logging enabled. SetLoggingPossiblyEnabled remains available
// for callers that need to enable the gate without producing a context, or
// to forcibly reset it (e.g. in tests). The runtime cost of the chain walk
// on every Lock/Unlock is significant enough that paying it only when at
// least one caller has indicated logging may be enabled is the right
// tradeoff: the historical default (deec5fb) was already "logging disabled
// by default", and no live caller in the dependency graph ever passes true.
var loggingPossiblyEnabled atomic.Bool

// deadlockDetectorPossiblyEnabled is the master switch for
// IsDeadlockDetectorEnabled. See loggingPossiblyEnabled for rationale,
// including the monotonic-on auto-flip behaviour wired through
// WithDeadlockDetectorEnabled / WithEnableDeadlock.
var deadlockDetectorPossiblyEnabled atomic.Bool

// SetLoggingPossiblyEnabled flips the global fast-path gate for
// IsLoggingEnabled. The gate auto-flips to true on the first
// WithLoggingEnabled(ctx, true) / WithNoLogging(ctx, false) call, so this
// setter is only needed for callers that want to enable logging without
// constructing a context, or to forcibly reset the gate (e.g. in tests).
// While the gate is false, per-context overrides via WithLoggingEnabled are
// ignored and IsLoggingEnabled always returns false without a context.Value
// chain walk.
func SetLoggingPossiblyEnabled(v bool) {
	loggingPossiblyEnabled.Store(v)
}

// SetDeadlockDetectorPossiblyEnabled flips the global fast-path gate for
// IsDeadlockDetectorEnabled. See SetLoggingPossiblyEnabled for rationale;
// the same auto-flip behaviour applies via WithDeadlockDetectorEnabled.
func SetDeadlockDetectorPossiblyEnabled(v bool) {
	deadlockDetectorPossiblyEnabled.Store(v)
}

type CtxKeyLogging struct{}

// WithNoLogging sets whether logging is disabled.
// As a side effect, calling it with noLogging == false flips the global
// loggingPossiblyEnabled gate on so the per-context override actually takes
// effect. The flip is monotonic-on; passing noLogging == true never flips
// the gate down.
//
// DEPRECATED: use WithLoggingEnabled
func WithNoLogging(ctx context.Context, noLogging bool) context.Context {
	if !noLogging {
		loggingPossiblyEnabled.Store(true)
	}
	return context.WithValue(ctx, CtxKeyLogging{}, !noLogging)
}

// IsNoLogging returns whether logging is disabled.
// DEPRECATED: use IsLoggingEnabled
func IsNoLogging(ctx context.Context) bool {
	if !loggingPossiblyEnabled.Load() {
		return true
	}
	v, _ := ctx.Value(CtxKeyLogging{}).(bool)
	return !v
}

// WithLoggingEnabled returns a context whose IsLoggingEnabled / IsNoLogging
// answer reflects the given logging value. As a side effect, calling it
// with logging == true flips the global loggingPossiblyEnabled gate on so
// the per-context override actually takes effect on the hot RWMutex path.
// The flip is monotonic-on; passing logging == false never flips the gate
// down because other live contexts may still want logging enabled.
func WithLoggingEnabled(ctx context.Context, logging bool) context.Context {
	if logging {
		loggingPossiblyEnabled.Store(true)
	}
	return context.WithValue(ctx, CtxKeyLogging{}, logging)
}

func IsLoggingEnabled(ctx context.Context) bool {
	if !loggingPossiblyEnabled.Load() {
		return false
	}
	v, _ := ctx.Value(CtxKeyLogging{}).(bool)
	return v
}

type CtxKeyEnableDeadlock struct{}

// WithEnableDeadlock sets whether the deadlock detector is enabled.
// Like WithDeadlockDetectorEnabled, calling it with enableDeadlock == true
// flips the global deadlockDetectorPossiblyEnabled gate on (monotonic-on).
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

// WithDeadlockDetectorEnabled returns a context whose
// IsDeadlockDetectorEnabled answer reflects the given enableDeadlock value.
// As a side effect, calling it with enableDeadlock == true flips the global
// deadlockDetectorPossiblyEnabled gate on so the per-context override
// actually takes effect on the hot RWMutex path. The flip is monotonic-on;
// passing enableDeadlock == false never flips the gate down because other
// live contexts may still want the detector enabled.
func WithDeadlockDetectorEnabled(ctx context.Context, enableDeadlock bool) context.Context {
	if enableDeadlock {
		deadlockDetectorPossiblyEnabled.Store(true)
	}
	return context.WithValue(ctx, CtxKeyEnableDeadlock{}, enableDeadlock)
}

func IsDeadlockDetectorEnabled(ctx context.Context) bool {
	if !deadlockDetectorPossiblyEnabled.Load() {
		return false
	}
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
