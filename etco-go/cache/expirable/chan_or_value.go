package expirable

import (
	"time"
)

type ChanOrValue[D any] struct {
	value *Expirable[D]
	chn   *ChanResult[D]
}

func NewCOVChan[D any](chn ChanResult[D]) ChanOrValue[D] {
	return ChanOrValue[D]{value: nil, chn: &chn}
}

func NewCOVValue[D any](value Expirable[D]) ChanOrValue[D] {
	return ChanOrValue[D]{value: &value, chn: nil}
}

func (cov ChanOrValue[D]) RecvExp() (D, time.Time, error) {
	if cov.value != nil {
		return cov.value.Data, cov.value.Expires, nil
	} else {
		return cov.chn.RecvExp()
	}
}

func (cov ChanOrValue[D]) RecvExpMin(prevExpCmp time.Time) (D, time.Time, error) {
	if cov.value != nil {
		if cov.value.Expires.After(prevExpCmp) {
			return cov.value.Data, prevExpCmp, nil
		} else {
			return cov.value.Data, cov.value.Expires, nil
		}
	} else {
		return cov.chn.RecvExpMin(prevExpCmp)
	}
}

func (cov ChanOrValue[D]) Recv() (Expirable[D], error) {
	if cov.value != nil {
		return *cov.value, nil
	} else {
		return cov.chn.Recv()
	}
}
