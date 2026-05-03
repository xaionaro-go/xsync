// context_test.go covers the fast-path gate behaviour for
// IsLoggingEnabled / IsDeadlockDetectorEnabled, including the
// monotonic-on auto-flip wired through the per-context setters.

package xsync

import (
	"context"
	"testing"
)

// resetGatesForTest forces both fast-path gates back to their pristine
// (false) state for the duration of one test. The gates are package-global
// atomic.Bool variables, so without an explicit reset a test that flips a
// gate on would bleed into every subsequent test.
func resetGatesForTest(t *testing.T) {
	t.Helper()
	loggingPossiblyEnabled.Store(false)
	deadlockDetectorPossiblyEnabled.Store(false)
	t.Cleanup(func() {
		loggingPossiblyEnabled.Store(false)
		deadlockDetectorPossiblyEnabled.Store(false)
	})
}

// --- logging gate -----------------------------------------------------------

func TestIsLoggingEnabled_DefaultsFalse(t *testing.T) {
	resetGatesForTest(t)

	if got := loggingPossiblyEnabled.Load(); got {
		t.Fatalf("loggingPossiblyEnabled: got %v, want false", got)
	}
	if got := IsLoggingEnabled(context.Background()); got {
		t.Fatalf("IsLoggingEnabled(Background): got %v, want false", got)
	}
}

func TestWithLoggingEnabledTrue_FlipsGate(t *testing.T) {
	resetGatesForTest(t)

	parent := context.Background()
	ctx := WithLoggingEnabled(parent, true)

	if got := loggingPossiblyEnabled.Load(); !got {
		t.Fatalf("loggingPossiblyEnabled after With...(true): got %v, want true", got)
	}
	if got := IsLoggingEnabled(ctx); !got {
		t.Fatalf("IsLoggingEnabled(returnedCtx): got %v, want true", got)
	}
	// Parent context must not have inherited the per-context override.
	if got := IsLoggingEnabled(parent); got {
		t.Fatalf("IsLoggingEnabled(parent): got %v, want false (parent unaffected)", got)
	}
}

func TestWithLoggingEnabledFalse_NoFlip(t *testing.T) {
	resetGatesForTest(t)

	ctx := WithLoggingEnabled(context.Background(), false)

	if got := loggingPossiblyEnabled.Load(); got {
		t.Fatalf("loggingPossiblyEnabled after With...(false): got %v, want false", got)
	}
	if got := IsLoggingEnabled(ctx); got {
		t.Fatalf("IsLoggingEnabled(ctx) with gate off: got %v, want false", got)
	}
}

func TestWithLoggingEnabled_MonotonicOn(t *testing.T) {
	resetGatesForTest(t)

	// Asking-true flips the gate.
	_ = WithLoggingEnabled(context.Background(), true)
	if !loggingPossiblyEnabled.Load() {
		t.Fatalf("gate should be on after asking-true")
	}
	// A subsequent asking-false call must NOT flip the gate down.
	_ = WithLoggingEnabled(context.Background(), false)
	if !loggingPossiblyEnabled.Load() {
		t.Fatalf("gate flipped down by asking-false; want monotonic-on")
	}
}

func TestWithNoLogging_AutoFlipsOnLoggingOn(t *testing.T) {
	resetGatesForTest(t)

	// WithNoLogging(ctx, false) means "logging is on" — must flip the gate.
	ctx := WithNoLogging(context.Background(), false)
	if !loggingPossiblyEnabled.Load() {
		t.Fatalf("WithNoLogging(false) did not flip gate on")
	}
	if !IsLoggingEnabled(ctx) {
		t.Fatalf("IsLoggingEnabled after WithNoLogging(false): want true")
	}
	if IsNoLogging(ctx) {
		t.Fatalf("IsNoLogging after WithNoLogging(false): want false")
	}
}

func TestWithNoLoggingTrue_NoFlip(t *testing.T) {
	resetGatesForTest(t)

	// WithNoLogging(ctx, true) means "logging off" — gate must stay down.
	_ = WithNoLogging(context.Background(), true)
	if loggingPossiblyEnabled.Load() {
		t.Fatalf("WithNoLogging(true) flipped gate on; want unchanged")
	}
}

func TestSetLoggingPossiblyEnabledFalse_DisablesFastPath(t *testing.T) {
	resetGatesForTest(t)

	ctx := WithLoggingEnabled(context.Background(), true)
	if !IsLoggingEnabled(ctx) {
		t.Fatalf("precondition: IsLoggingEnabled should be true after asking-true")
	}

	// Force the gate back off — the per-context override stays in the
	// context, but the predicate must short-circuit to false.
	SetLoggingPossiblyEnabled(false)
	if IsLoggingEnabled(ctx) {
		t.Fatalf("IsLoggingEnabled after gate forced off: want false")
	}

	// Re-enabling the gate (without a new context) restores visibility of
	// the still-attached per-context true.
	SetLoggingPossiblyEnabled(true)
	if !IsLoggingEnabled(ctx) {
		t.Fatalf("IsLoggingEnabled after gate forced back on: want true")
	}
}

// --- deadlock-detector gate -------------------------------------------------

func TestIsDeadlockDetectorEnabled_DefaultsFalse(t *testing.T) {
	resetGatesForTest(t)

	if got := deadlockDetectorPossiblyEnabled.Load(); got {
		t.Fatalf("deadlockDetectorPossiblyEnabled: got %v, want false", got)
	}
	if got := IsDeadlockDetectorEnabled(context.Background()); got {
		t.Fatalf("IsDeadlockDetectorEnabled(Background): got %v, want false", got)
	}
}

func TestWithDeadlockDetectorEnabledTrue_FlipsGate(t *testing.T) {
	resetGatesForTest(t)

	parent := context.Background()
	ctx := WithDeadlockDetectorEnabled(parent, true)

	if !deadlockDetectorPossiblyEnabled.Load() {
		t.Fatalf("deadlockDetectorPossiblyEnabled after With...(true): want true")
	}
	if !IsDeadlockDetectorEnabled(ctx) {
		t.Fatalf("IsDeadlockDetectorEnabled(returnedCtx): want true")
	}
	if IsDeadlockDetectorEnabled(parent) {
		t.Fatalf("IsDeadlockDetectorEnabled(parent): want false (parent unaffected)")
	}
}

func TestWithDeadlockDetectorEnabledFalse_NoFlip(t *testing.T) {
	resetGatesForTest(t)

	ctx := WithDeadlockDetectorEnabled(context.Background(), false)

	if deadlockDetectorPossiblyEnabled.Load() {
		t.Fatalf("deadlockDetectorPossiblyEnabled after With...(false): want false")
	}
	if IsDeadlockDetectorEnabled(ctx) {
		t.Fatalf("IsDeadlockDetectorEnabled(ctx) with gate off: want false")
	}
}

func TestWithDeadlockDetectorEnabled_MonotonicOn(t *testing.T) {
	resetGatesForTest(t)

	_ = WithDeadlockDetectorEnabled(context.Background(), true)
	if !deadlockDetectorPossiblyEnabled.Load() {
		t.Fatalf("gate should be on after asking-true")
	}
	_ = WithDeadlockDetectorEnabled(context.Background(), false)
	if !deadlockDetectorPossiblyEnabled.Load() {
		t.Fatalf("gate flipped down by asking-false; want monotonic-on")
	}
}

func TestWithEnableDeadlockTrue_FlipsGate(t *testing.T) {
	resetGatesForTest(t)

	// WithEnableDeadlock is the deprecated alias; verify it inherits the
	// auto-flip behaviour (since it just delegates).
	ctx := WithEnableDeadlock(context.Background(), true)
	if !deadlockDetectorPossiblyEnabled.Load() {
		t.Fatalf("WithEnableDeadlock(true) did not flip gate on")
	}
	if !IsDeadlockDetectorEnabled(ctx) {
		t.Fatalf("IsDeadlockDetectorEnabled after WithEnableDeadlock(true): want true")
	}
}

func TestSetDeadlockDetectorPossiblyEnabledFalse_DisablesFastPath(t *testing.T) {
	resetGatesForTest(t)

	ctx := WithDeadlockDetectorEnabled(context.Background(), true)
	if !IsDeadlockDetectorEnabled(ctx) {
		t.Fatalf("precondition: IsDeadlockDetectorEnabled should be true after asking-true")
	}

	SetDeadlockDetectorPossiblyEnabled(false)
	if IsDeadlockDetectorEnabled(ctx) {
		t.Fatalf("IsDeadlockDetectorEnabled after gate forced off: want false")
	}

	SetDeadlockDetectorPossiblyEnabled(true)
	if !IsDeadlockDetectorEnabled(ctx) {
		t.Fatalf("IsDeadlockDetectorEnabled after gate forced back on: want true")
	}
}

// --- gate independence ------------------------------------------------------

func TestGates_AreIndependent(t *testing.T) {
	resetGatesForTest(t)

	_ = WithLoggingEnabled(context.Background(), true)
	if !loggingPossiblyEnabled.Load() {
		t.Fatalf("logging gate should be on")
	}
	if deadlockDetectorPossiblyEnabled.Load() {
		t.Fatalf("deadlock-detector gate must not flip when only logging asked-true")
	}

	// Reset between sub-checks; cleanup already registered above.
	loggingPossiblyEnabled.Store(false)
	deadlockDetectorPossiblyEnabled.Store(false)

	_ = WithDeadlockDetectorEnabled(context.Background(), true)
	if !deadlockDetectorPossiblyEnabled.Load() {
		t.Fatalf("deadlock-detector gate should be on")
	}
	if loggingPossiblyEnabled.Load() {
		t.Fatalf("logging gate must not flip when only deadlock-detector asked-true")
	}
}
