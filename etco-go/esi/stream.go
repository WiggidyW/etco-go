package esi

import (
	"github.com/WiggidyW/etco-go/cache/expirable"
)

type RepOrStream[E any] struct {
	Rep    *[]E
	Stream *expirable.ChanResult[[]E]
}

func newBootstrapRepOrStream[E any]() RepOrStream[E] {
	return RepOrStream[E]{Rep: new([]E), Stream: nil}
}
