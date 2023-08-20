package util

import "context"

func SendUntilDone[T any](
	ctx context.Context,
	chn chan<- T,
	t T,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case chn <- t:
		return nil
	}
}

type ChanResult[T any] struct {
	ctx    context.Context
	chnOk  chan T
	chnErr chan error
}

func NewChanResult[T any](ctx context.Context) ChanResult[T] {
	return ChanResult[T]{
		ctx:    ctx,
		chnOk:  make(chan T),
		chnErr: make(chan error),
	}
}

func (rc ChanResult[T]) SendOk(t T) error {
	return SendUntilDone[T](rc.ctx, rc.chnOk, t)
}

func (rc ChanResult[T]) SendErr(err error) error {
	return SendUntilDone[error](rc.ctx, rc.chnErr, err)
}

func (rc ChanResult[T]) Wait() error {
	<-rc.ctx.Done()
	return rc.ctx.Err()
}

func (rc ChanResult[T]) Recv() (t T, err error) {
	select {
	case t = <-rc.chnOk:
		return t, nil
	case err = <-rc.chnErr:
		return t, err
	case <-rc.ctx.Done():
		return t, rc.ctx.Err()
	}
}

func (rc ChanResult[T]) RecvAll(len int) ([]T, error) {
	ts := make([]T, 0, len)
	for i := 0; i < len; i++ {
		t, err := rc.Recv()
		if err != nil {
			return nil, err
		} else {
			ts = append(ts, t)
		}
	}
	return ts, nil
}

func (rc ChanResult[T]) RecvAllDiscard(len int) error {
	for i := 0; i < len; i++ {
		_, err := rc.Recv()
		if err != nil {
			return err
		}
	}
	return nil
}

func (rc ChanResult[T]) Split() (ChanSendResult[T], ChanRecvResult[T]) {
	return rc.ToSend(), rc.ToRecv()
}

func (rc ChanResult[T]) ToSend() ChanSendResult[T] {
	return ChanSendResult[T]{inner: rc}
}

func (rc ChanResult[T]) ToRecv() ChanRecvResult[T] {
	return ChanRecvResult[T]{inner: rc}
}

type ChanSendResult[T any] struct{ inner ChanResult[T] }

func (src ChanSendResult[T]) SendOk(t T) error        { return src.inner.SendOk(t) }
func (src ChanSendResult[T]) SendErr(err error) error { return src.inner.SendErr(err) }
func (src ChanSendResult[T]) Wait() error             { return src.inner.Wait() }

type ChanRecvResult[T any] struct{ inner ChanResult[T] }

func (rrc ChanRecvResult[T]) Recv() (T, error)             { return rrc.inner.Recv() }
func (rrc ChanRecvResult[T]) RecvAll(len int) ([]T, error) { return rrc.inner.RecvAll(len) }
func (rrc ChanRecvResult[T]) RecvAllDiscard(len int) error { return rrc.inner.RecvAllDiscard(len) }
func (rrc ChanRecvResult[T]) Wait() error                  { return rrc.inner.Wait() }
