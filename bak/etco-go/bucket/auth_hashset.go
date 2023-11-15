package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	AUTH_HASH_SET_EXPIRES_IN time.Duration = 24 * time.Hour
	AUTH_HASH_SET_BUF_CAP    int           = 0
)

func init() {
	keys.TypeStrAuthHashSet = cache.RegisterType[b.AuthHashSet]("authhashset", AUTH_HASH_SET_BUF_CAP)
}

func GetAuthHashSet(
	x cache.Context,
	domain string,
) (
	rep b.AuthHashSet,
	expires time.Time,
	err error,
) {
	return get(
		x,
		func(ctx context.Context) (b.AuthHashSet, error) {
			return client.ReadAuthHashSet(ctx, domain)
		},
		keys.CacheKeyAuthHashSet(domain),
		keys.TypeStrAuthHashSet,
		AUTH_HASH_SET_EXPIRES_IN,
		nil,
	)
}

func SetAuthHashSet(
	x cache.Context,
	domain string,
	rep b.AuthHashSet,
) (
	err error,
) {
	return set(
		x,
		func(ctx context.Context, rep b.AuthHashSet) error {
			return client.WriteAuthHashSet(ctx, rep, domain)
		},
		keys.CacheKeyAuthHashSet(domain),
		keys.TypeStrAuthHashSet,
		AUTH_HASH_SET_EXPIRES_IN,
		rep,
		nil,
	)
}
