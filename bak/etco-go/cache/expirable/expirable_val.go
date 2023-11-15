package expirable

import (
	"time"
)

type ExpirableVal[D any] struct {
	Data    D
	Expires time.Time
}

func (e ExpirableVal[any]) Expired() bool {
	return time.Now().After(e.Expires)
}

func NewVal[D any](data D, expires time.Time) ExpirableVal[D] {
	return ExpirableVal[D]{Data: data, Expires: expires}
}

func NewValPtr[D any](data D, expires time.Time) *ExpirableVal[D] {
	expirable := NewVal(data, expires)
	return &expirable
}
