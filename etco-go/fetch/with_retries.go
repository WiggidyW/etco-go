package fetch

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
)

func FetchWithRetries[REP any](
	x cache.Context,
	fetch Fetch[REP],
	numRetries int,
	shouldRetry func(error) bool,
) (
	rep REP,
	expires time.Time,
	err error,
) {
	return fetchWithRetries(x, fetch, numRetries, shouldRetry, 0)
}

func fetchWithRetries[REP any](
	x cache.Context,
	fetch Fetch[REP],
	numRetries int,
	shouldRetry func(error) bool,
	attempt int,
) (
	rep REP,
	expires time.Time,
	err error,
) {
	rep, expires, err = fetch(x)
	if err != nil && attempt < numRetries && shouldRetry(err) {
		return fetchWithRetries(
			x,
			fetch,
			numRetries,
			shouldRetry,
			attempt+1,
		)
	} else {
		return rep, expires, err
	}
}
