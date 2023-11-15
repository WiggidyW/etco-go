package servercache

import (
	"fmt"
	"time"
)

type ErrServerUnlock struct {
	err error
}

func (e ErrServerUnlock) Error() string {
	return "ErrServerUnlock: " + e.err.Error()
}

type ErrServerObtainLock struct {
	err error
}

func (e ErrServerObtainLock) Error() string {
	return "ErrServerLock: " + e.err.Error()
}

type ErrServerRefreshLock struct {
	err error
}

func (e ErrServerRefreshLock) Error() string {
	return "ErrServerRefreshLock: " + e.err.Error()
}

type ErrServerGet struct {
	err error
}

func (e ErrServerGet) Error() string {
	return "ErrServerGet: " + e.err.Error()
}

type ErrServerSet struct {
	err error
}

func (e ErrServerSet) Error() string {
	return "ErrServerSet: " + e.err.Error()
}

type ErrServerDel struct {
	err error
}

func (e ErrServerDel) Error() string {
	return "ErrServerDel: " + e.err.Error()
}

type ErrInvalidSet struct {
	key string
	ttl time.Duration
}

func (e ErrInvalidSet) Error() string {
	return fmt.Sprintf(
		"ErrInvalidSet: cannot set expired value (key: %s, ttl: %s)",
		e.key,
		e.ttl,
	)
}
