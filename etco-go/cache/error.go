package cache

import (
	"fmt"

	"github.com/WiggidyW/etco-go/cache/keys"
)

type CacheLockErr struct {
	err   error
	Key   keys.Key
	Scope int64
}

func (e CacheLockErr) Unwrap() error { return e.err }
func (e CacheLockErr) Error() string {
	return fmt.Sprintf(
		"LockErr: key: '%s', scope: '%d', err: %s",
		e.Key.PrettyString(),
		e.Scope,
		e.err.Error(),
	)
}

type LockNil struct{}

func (LockNil) Unwrap() error { return nil }
func (LockNil) Error() string { return "lock is nil" }

type CacheErr struct {
	err    error
	Key    keys.Key
	Scope  int64
	Method string
}

func (e CacheErr) Unwrap() error { return e.err }
func (e CacheErr) Error() string {
	return fmt.Sprintf(
		"CacheErr: key: '%s', scope: '%d', method: '%s', err: %s",
		e.Key.PrettyString(),
		e.Scope,
		e.Method,
		e.err.Error(),
	)
}
