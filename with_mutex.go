// with_mutex.go defines a utility struct that pairs a generic value with an RWMutex.

package xsync

type WithMutex[T any] struct {
	RWMutex
	Value T
}
