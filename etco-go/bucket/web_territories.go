package bucket

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/proto/protoerr"
)

func protoMergeSetTerritories[U any, T any](
	x cache.Context,
	updates U,
	getBundleKeys func(cache.Context) (map[string]struct{}, time.Time, error),
	getOriginal func(cache.Context) (T, time.Time, error),
	mergeUpdates func(T, U, map[string]struct{}) error,
	setUpdated func(cache.Context, T) error,
) (
	err error,
) {
	// fetch bundle keys in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chnBundleKeys := expirable.NewChanResult[map[string]struct{}](x.Ctx(), 1, 0)
	go expirable.P1Transceive(chnBundleKeys, x, getBundleKeys)

	// fetch the original territories
	var territories T
	territories, _, err = getOriginal(x)
	if err != nil {
		return err
	}

	// recv bundle keys
	var bundleKeys map[string]struct{}
	bundleKeys, _, err = chnBundleKeys.RecvExp()
	if err != nil {
		return err
	}

	// merge updates
	err = mergeUpdates(territories, updates, bundleKeys)
	if err != nil {
		return protoerr.New(protoerr.INVALID_MERGE, err)
	}

	// set updated
	return setUpdated(x, territories)
}
