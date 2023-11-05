package cache

type SetLocalOnServerHit[REP any] interface {
	SetLocalOnServerHit(rep *REP) bool
}

func NewSloshFunc[REP any](f func(rep *REP) bool) *SetLocalOnServerHit[REP] {
	sloshFunc := sloshFunc[REP](f)
	slosh := SetLocalOnServerHit[REP](sloshFunc)
	return &slosh
}

func NewSloshBool[REP any](b bool) *SetLocalOnServerHit[REP] {
	sloshBool := sloshBool[REP](b)
	slosh := SetLocalOnServerHit[REP](sloshBool)
	return &slosh
}

type sloshFunc[REP any] func(rep *REP) bool

func (f sloshFunc[REP]) SetLocalOnServerHit(rep *REP) bool {
	return f(rep)
}

type sloshBool[REP any] bool

func (b sloshBool[REP]) SetLocalOnServerHit(rep *REP) bool {
	return bool(b)
}

func setLocalOnServerHitOrDefault[REP any](
	slosh *SetLocalOnServerHit[REP],
	def bool,
) bool {
	if slosh == nil {
		return def
	}
	return (*slosh).SetLocalOnServerHit(nil)
}
