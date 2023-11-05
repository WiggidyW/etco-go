package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	AUTH_HASH_SET_BUF_CAP          int           = 0
	AUTH_HASH_SET_LOCK_TTL         time.Duration = 1 * time.Minute
	AUTH_HASH_SET_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	AUTH_HASH_SET_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrAuthHashSet = localcache.RegisterType[b.AuthHashSet](AUTH_HASH_SET_BUF_CAP)
}

func GetAuthHashSet(
	ctx context.Context,
	domain string,
) (
	rep b.AuthHashSet,
	expires *time.Time,
	err error,
) {
	return get(
		ctx,
		func(ctx context.Context) (b.AuthHashSet, error) {
			return client.ReadAuthHashSet(ctx, domain)
		},
		keys.TypeStrAuthHashSet,
		keys.CacheKeyAuthHashSet(domain),
		AUTH_HASH_SET_LOCK_TTL,
		AUTH_HASH_SET_LOCK_MAX_BACKOFF,
		AUTH_HASH_SET_EXPIRES_IN,
		nil,
	)
}

func SetAuthHashSet(
	ctx context.Context,
	domain string,
	rep b.AuthHashSet,
) (
	err error,
) {
	return set(
		ctx,
		func(ctx context.Context, rep b.AuthHashSet) error {
			return client.WriteAuthHashSet(ctx, rep, domain)
		},
		keys.TypeStrAuthHashSet,
		keys.CacheKeyAuthHashSet(domain),
		AUTH_HASH_SET_LOCK_TTL,
		AUTH_HASH_SET_LOCK_MAX_BACKOFF,
		AUTH_HASH_SET_EXPIRES_IN,
		rep,
	)
}
