package cache

import "time"

type Expirable[D any] interface {
	Data() D
	Expires() time.Time
	Cache() bool
}

type ExpirableData[D any] struct {
	Ddata    D
	Eexpires time.Time
}

func NewExpirableData[D any](data D, expiry time.Time) ExpirableData[D] {
	return ExpirableData[D]{data, expiry}
}

func NewExpirableDataPtr[D any](data D, expiry time.Time) *ExpirableData[D] {
	return &ExpirableData[D]{data, expiry}
}

func (ed ExpirableData[D]) Expires() time.Time {
	return ed.Eexpires
}

func (ed ExpirableData[D]) Data() D {
	return ed.Ddata
}

func (ExpirableData[D]) Cache() bool {
	return true
}

type MaybeCacheExpirableData[D any] struct {
	Ddata    D
	Eexpires time.Time
	Ccache   bool
}

func NewMaybeCacheExpirableData[D any](data D, expiry time.Time, cache bool) MaybeCacheExpirableData[D] {
	return MaybeCacheExpirableData[D]{data, expiry, cache}
}

func NewMaybeCacheExpirableDataPtr[D any](data D, expiry time.Time, cache bool) *MaybeCacheExpirableData[D] {
	return &MaybeCacheExpirableData[D]{data, expiry, cache}
}

func (ed MaybeCacheExpirableData[D]) Expires() time.Time {
	return ed.Eexpires
}

func (ed MaybeCacheExpirableData[D]) Data() D {
	return ed.Ddata
}

func (ed MaybeCacheExpirableData[D]) Cache() bool {
	return ed.Ccache
}
