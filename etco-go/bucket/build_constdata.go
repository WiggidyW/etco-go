package bucket

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	BUILD_CONST_DATA_BUF_CAP    int           = 0
	BUILD_CONST_DATA_EXPIRES_IN time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrBuildConstData = cache.RegisterType[b.ConstantsData]("buildconstdata", BUILD_CONST_DATA_BUF_CAP)
}

func GetBuildConstData(
	x cache.Context,
) (
	rep b.ConstantsData,
	expires time.Time,
	err error,
) {
	return get(
		x,
		client.ReadConstantsData,
		keys.CacheKeyBuildConstData,
		keys.TypeStrBuildConstData,
		BUILD_CONST_DATA_EXPIRES_IN,
		nil,
	)
}

func SetBuildConstData(
	x cache.Context,
	rep b.ConstantsData,
) (
	err error,
) {
	return set(
		x,
		client.WriteConstantsData,
		keys.CacheKeyBuildConstData,
		keys.TypeStrBuildConstData,
		BUILD_CONST_DATA_EXPIRES_IN,
		rep,
		nil,
	)
}
