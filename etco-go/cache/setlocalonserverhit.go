package cache

type SetLocalOnServerHit[REP any] func(rep *REP) bool

func SloshTrue[REP any](rep *REP) bool  { return true }
func SloshFalse[REP any](rep *REP) bool { return false }
