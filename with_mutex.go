package xsync

type WithMutex[T any] struct {
	RWMutex
	Value T
}
