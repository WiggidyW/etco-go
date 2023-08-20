package naiveclient

import (
	"context"
	"time"

	"github.com/WiggidyW/weve-esi/util"
)

type ChanPage[M any] struct {
	inner       util.ChanResult[M]
	pages       int32
	headExpires time.Time
}

func NewChanPage[M any](
	ctx context.Context,
	Pages int32,
	HeadExpires time.Time,
) ChanPage[M] {
	return ChanPage[M]{
		inner:       util.NewChanResult[M](ctx),
		pages:       Pages,
		headExpires: HeadExpires,
	}
}

func (rc ChanPage[M]) NumPages() int32 {
	return rc.pages
}

func (rc ChanPage[M]) HeadExpires() time.Time {
	return rc.headExpires
}

func (rc ChanPage[M]) SendOk(m M) error {
	return rc.inner.SendOk(m)
}

func (rc ChanPage[M]) SendErr(err error) error {
	return rc.inner.SendErr(err)
}

func (rc ChanPage[M]) Recv() (M, error) {
	return rc.inner.Recv()
}

func (rc ChanPage[M]) RecvAll() ([]M, error) {
	return rc.inner.RecvAll(int(rc.pages))
}

func (rc ChanPage[M]) Split() (ChanSendPage[M], ChanRecvPage[M]) {
	return rc.ToSend(), rc.ToRecv()
}

func (rc ChanPage[M]) ToSend() ChanSendPage[M] {
	return ChanSendPage[M]{inner: rc}
}

func (rc ChanPage[M]) ToRecv() ChanRecvPage[M] {
	return ChanRecvPage[M]{inner: rc}
}

type ChanSendPage[M any] struct{ inner ChanPage[M] }

func (rc ChanSendPage[M]) NumPages() int32         { return rc.inner.NumPages() }
func (rc ChanSendPage[M]) HeadExpires() time.Time  { return rc.inner.HeadExpires() }
func (rc ChanSendPage[M]) SendOk(m M) error        { return rc.inner.SendOk(m) }
func (rc ChanSendPage[M]) SendErr(err error) error { return rc.inner.SendErr(err) }

type ChanRecvPage[M any] struct{ inner ChanPage[M] }

func (rc ChanRecvPage[M]) NumPages() int32        { return rc.inner.NumPages() }
func (rc ChanRecvPage[M]) HeadExpires() time.Time { return rc.inner.HeadExpires() }
func (rc ChanRecvPage[M]) Recv() (M, error)       { return rc.inner.Recv() }
func (rc ChanRecvPage[M]) RecvAll() ([]M, error)  { return rc.inner.RecvAll() }
