package pbmerge

import (
	cfgerr "github.com/WiggidyW/eve-trading-co-go/error/configerror"
)

func newError(name string, errStr string) cfgerr.ErrInvalid {
	return cfgerr.ErrInvalid{Err: cfgerr.ErrMarketInvalid{
		Market:    name,
		ErrString: errStr,
	}}
}
