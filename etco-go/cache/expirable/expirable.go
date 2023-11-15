package expirable

import (
	"time"
)

type Expirable[D any] struct {
	Data    D
	Expires time.Time
}

func (e Expirable[any]) Expired() bool {
	return time.Now().After(e.Expires)
}

func New[D any](data D, expires time.Time) Expirable[D] {
	return Expirable[D]{Data: data, Expires: expires}
}

func NewPtr[D any](data D, expires time.Time) *Expirable[D] {
	expirable := New(data, expires)
	return &expirable
}

func NewMarshal[D any](data D) Expirable[D] {
	return Expirable[D]{Data: data}
}

func NewMarshalPtr[D any](data D) *Expirable[D] {
	expirable := NewMarshal(data)
	return &expirable
}

func NewAny(data any, expires time.Time) Expirable[any] {
	return Expirable[any]{Data: &data, Expires: expires}
}
