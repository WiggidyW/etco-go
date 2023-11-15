package util

import "github.com/WiggidyW/chanresult"

func TransceiveDone(
	chn chanresult.ChanResult[struct{}],
	fn func(),
) error {
	fn()
	return chn.SendOk(struct{}{})
}

func TransceiveNoErr[T any](
	chn chanresult.ChanResult[T],
	fn func() T,
) error {
	return chn.SendOk(fn())
}

func Transceive[T any](
	chn chanresult.ChanResult[T],
	fn func() (T, error),
) error {
	t, err := fn()
	if err != nil {
		return chn.SendErr(err)
	} else {
		return chn.SendOk(t)
	}
}
