package cache

import "time"

type Expirable[D any] interface {
	Data() D
	Expires() time.Time
}

type ExpirableData[D any] struct {
	data    D
	expires time.Time
}

func NewExpirableData[D any](data D, expiry time.Time) ExpirableData[D] {
	return ExpirableData[D]{data, expiry}
}

func (ed ExpirableData[D]) Expires() time.Time {
	return ed.expires
}

func (ed ExpirableData[D]) Data() D {
	return ed.data
}
