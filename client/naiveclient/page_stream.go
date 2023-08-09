package naiveclient

import "time"

type PageStream[M any] struct {
	okChn      chan *M
	errChn     chan error
	pages      int32
	headExpiry time.Time
}

func makePageStream[M any](pages int32, headExpires time.Time) PageStream[M] {
	return PageStream[M]{
		okChn:      make(chan *M, pages),
		errChn:     make(chan error, pages),
		pages:      pages,
		headExpiry: headExpires,
	}
}

func (ps *PageStream[M]) Recv() (*M, error) {
	select {
	case ok := <-ps.okChn:
		return ok, nil
	case err := <-ps.errChn:
		return nil, err
	}
}

func (pc *PageStream[M]) Close() {
	close(pc.okChn)
	close(pc.errChn)
}

func (pc *PageStream[M]) NumPages() int32 {
	return pc.pages
}

func (pc *PageStream[M]) HeadExpires() time.Time {
	return pc.headExpiry
}

func (ps *PageStream[M]) sendOk(ok *M) {
	ps.okChn <- ok
}

func (ps *PageStream[M]) sendErr(err error) {
	ps.errChn <- err
}
