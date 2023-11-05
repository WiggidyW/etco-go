package fetch

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func HandleFetch[REP any](
	ctx context.Context,
	preFetchParams *prefetch.Params[REP],
	fetch Fetch[REP],
) (
	rep *REP,
	expires *time.Time,
	err error,
) {
	var preFetchData *prefetch.UnhandledData
	if preFetchParams != nil {
		var expirableRep *expirable.Expirable[REP]
		expirableRep, preFetchData, err = prefetch.Handle(
			ctx,
			*preFetchParams,
		)
		if err != nil {
			return nil, nil, err
		} else if expirableRep != nil {
			rep, expires = expirableRep.Data, expirableRep.Expires
			return rep, expires, nil
		}
	}

	var postFetchParams *postfetch.Params
	rep, expires, postFetchParams, err = fetch(ctx)
	go postfetch.Handle(preFetchData, rep, expires, err, postFetchParams)

	return rep, expires, err
}
