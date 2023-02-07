package executioncontext

import "context"

// Execution context flags
const (
	LIVE   Flag = 1 << iota // live events handling
	REPLAY                  // replay events handling
)

// Flag type
type Flag uint8

// set flag
func (f Flag) set(flag Flag) Flag { return f | flag }

// clear flag
func (f Flag) clear(flag Flag) Flag { return f &^ flag }

// toggle flag
func (f Flag) toggle(flag Flag) Flag { return f ^ flag }

// has flag
func (f Flag) has(flag Flag) bool { return f&flag != 0 }

type key struct{}

// WithFlag returns a new Context that carries flag.
func WithFlag(ctx context.Context, flag Flag) context.Context {
	if ctx == nil {
		return nil
	}

	flags, ok := ctx.Value(key{}).(Flag)
	if !ok {
		flags = 0
	}

	return context.WithValue(ctx, key{}, flags.set(flag))
}

// ClearFlag returns a new Context that no longer carries flag.
func ClearFlag(ctx context.Context, flag Flag) context.Context {
	if ctx == nil {
		return nil
	}

	flags, ok := ctx.Value(key{}).(Flag)
	if !ok {
		flags = 0
	}

	return context.WithValue(ctx, key{}, flags.clear(flag))
}

// ToggleFlag returns a new Context with toggled flag.
func ToggleFlag(ctx context.Context, flag Flag) context.Context {
	if ctx == nil {
		return nil
	}

	flags, ok := ctx.Value(key{}).(Flag)
	if !ok {
		flags = 0
	}

	return context.WithValue(ctx, key{}, flags.toggle(flag))
}

// FromContext returns the slice of flags stored in ctx
func FromContext(ctx context.Context) Flag {
	if ctx == nil {
		return 0
	}

	flags, ok := ctx.Value(key{}).(Flag)
	if !ok {
		return 0
	}

	return flags
}

// has returns the bool value based on flag occurrence in context.
func Has(ctx context.Context, flag Flag) bool {
	flags, ok := ctx.Value(key{}).(Flag)
	if !ok {
		return false
	}

	return flags.has(flag)
}
