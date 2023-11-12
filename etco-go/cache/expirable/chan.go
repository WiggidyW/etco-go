package expirable

import (
	"context"
	"time"

	"github.com/WiggidyW/chanresult"
)

type ChanResult[D any] struct {
	chanresult.ChanResult[ExpirableVal[D]]
}

func NewChanResult[D any](
	ctx context.Context,
	okCap, errCap int,
) ChanResult[D] {
	return ChanResult[D]{
		chanresult.NewChanResult[ExpirableVal[D]](ctx, okCap, errCap),
	}
}

func (chn ChanResult[D]) SendExp(
	data D,
	expires time.Time,
	err error,
) error {
	if err != nil {
		return chn.SendErr(err)
	} else {
		return chn.SendExpOk(data, expires)
	}
}

func (chn ChanResult[D]) SendExpOk(
	data D,
	expires time.Time,
) error {
	return chn.SendOk(NewVal(data, expires))
}

func (chn ChanResult[D]) RecvExp() (D, time.Time, error) {
	expirable, err := chn.Recv()
	if err != nil {
		return *new(D), time.Time{}, err
	} else {
		return expirable.Data, expirable.Expires, nil
	}
}

func (chn ChanResult[D]) RecvExpMin(prevExpCmp time.Time) (D, time.Time, error) {
	expirable, err := chn.Recv()
	if err != nil {
		return *new(D), prevExpCmp, err
	} else if expirable.Expires.After(prevExpCmp) {
		return expirable.Data, prevExpCmp, nil
	} else {
		return expirable.Data, expirable.Expires, nil
	}
}

func (chn ChanResult[D]) Transceive(
	fn func() (D, time.Time, error),
) error {
	return chn.SendExp(fn())
}

func Param1Transceive[D any, P1 any](
	chn ChanResult[D],
	p1 P1,
	fn func(P1) (D, time.Time, error),
) error {
	return chn.SendExp(fn(p1))
}
func Param2Transceive[D any, P1 any, P2 any](
	chn ChanResult[D],
	p1 P1,
	p2 P2,
	fn func(P1, P2) (D, time.Time, error),
) error {
	return chn.SendExp(fn(p1, p2))
}
func Param3Transceive[D any, P1 any, P2 any, P3 any](
	chn ChanResult[D],
	p1 P1,
	p2 P2,
	p3 P3,
	fn func(P1, P2, P3) (D, time.Time, error),
) error {
	return chn.SendExp(fn(p1, p2, p3))
}
func Param4Transceive[D any, P1 any, P2 any, P3 any, P4 any](
	chn ChanResult[D],
	p1 P1,
	p2 P2,
	p3 P3,
	p4 P4,
	fn func(P1, P2, P3, P4) (D, time.Time, error),
) error {
	return chn.SendExp(fn(p1, p2, p3, p4))
}
