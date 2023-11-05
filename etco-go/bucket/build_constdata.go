package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	BUILD_CONST_DATA_BUF_CAP          int           = 0
	BUILD_CONST_DATA_LOCK_TTL         time.Duration = 1 * time.Minute
	BUILD_CONST_DATA_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	BUILD_CONST_DATA_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrBuildConstData = localcache.RegisterType[b.ConstantsData](BUILD_CONST_DATA_BUF_CAP)
}

func GetBuildConstData(
	ctx context.Context,
) (
	rep b.ConstantsData,
	expires *time.Time,
	err error,
) {
	return get(
		ctx,
		client.ReadConstantsData,
		keys.TypeStrBuildConstData,
		keys.CacheKeyBuildConstData,
		BUILD_CONST_DATA_LOCK_TTL,
		BUILD_CONST_DATA_LOCK_MAX_BACKOFF,
		BUILD_CONST_DATA_EXPIRES_IN,
		nil,
	)
}

func SetBuildConstData(
	ctx context.Context,
	rep b.ConstantsData,
) (
	err error,
) {
	return set(
		ctx,
		client.WriteConstantsData,
		keys.TypeStrBuildConstData,
		keys.CacheKeyBuildConstData,
		BUILD_CONST_DATA_LOCK_TTL,
		BUILD_CONST_DATA_LOCK_MAX_BACKOFF,
		BUILD_CONST_DATA_EXPIRES_IN,
		rep,
	)
}
