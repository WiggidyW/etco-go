package cache

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

type ErrLocalUnlock struct {
	err error
}

func (e ErrLocalUnlock) Error() string {
	return "ErrLocalUnlock: " + e.err.Error()
}

type ErrServerLock struct {
	err error
}

func (e ErrServerLock) Error() string {
	return "ErrServerLock: " + e.err.Error()
}

type ErrServerGet struct {
	err error
}

func (e ErrServerGet) Error() string {
	return "ErrServerGet: " + e.err.Error()
}

type ErrLocalDeserialize struct {
	err error
}

func (e ErrLocalDeserialize) Error() string {
	return "ErrLocalDeserialize: " + e.err.Error()
}

type ErrServerDeserialize struct {
	err error
}

func (e ErrServerDeserialize) Error() string {
	return "ErrServerDeserialize: " + e.err.Error()
}

type ErrSerialize struct {
	err error
}

func (e ErrSerialize) Error() string {
	return "ErrSerialize: " + e.err.Error()
}

type ErrServerSet struct {
	err error
}

func (e ErrServerSet) Error() string {
	return "ErrServerSet: " + e.err.Error()
}

type ErrInvalidLock struct {
	funcName string
}

func (e ErrInvalidLock) Error() string {
	return fmt.Sprintf(
		"ErrInvalidLock: %s called with invalid lock",
		e.funcName,
	)
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
