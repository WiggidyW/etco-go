package fetch

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/fetch/postfetch"
)

type HandledFetchVal[REP any] func(context.Context) (
	REP,
	*time.Time,
	error,
)

type HandledFetch[REP any] func(context.Context) (
	*REP,
	*time.Time,
	error,
)

type FetchVal[REP any] func(context.Context) (
	REP,
	*time.Time,
	*postfetch.Params,
	error,
)

type Fetch[REP any] func(context.Context) (
	*REP,
	*time.Time,
	*postfetch.Params,
	error,
)
