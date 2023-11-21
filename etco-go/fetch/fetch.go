package fetch

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
)

type HandledFetch[REP any] func(cache.Context) (
	rep REP,
	expires time.Time,
	err error,
)

type Fetch[REP any] func(cache.Context) (
	rep REP,
	expires time.Time,
	err error,
)

type CachingFetch[REP any] func(cache.Context) (
	rep REP,
	expires time.Time,
	postFetch *cachepostfetch.Params,
	err error,
)
