package bucket

import (
	"fmt"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_SHOP_LOCATIONS_BUF_CAP    int           = 0
	WEB_SHOP_LOCATIONS_EXPIRES_IN time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebShopLocations = cache.RegisterType[map[b.LocationId]b.WebShopLocation]("webshoplocations", WEB_SHOP_LOCATIONS_BUF_CAP)
}

func GetWebShopLocations(
	x cache.Context,
) (
	rep map[b.LocationId]b.WebShopLocation,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebShopLocations,
		keys.CacheKeyWebShopLocations,
		keys.TypeStrWebShopLocations,
		WEB_SHOP_LOCATIONS_EXPIRES_IN,
		build.CAPACITY_WEB_SHOP_LOCATIONS,
	)
}

func ProtoGetWebShopLocations(
	x cache.Context,
) (
	rep map[int64]*proto.CfgShopLocation,
	expires time.Time,
	err error,
) {
	var webShopLocations map[b.LocationId]b.WebShopLocation
	webShopLocations, expires, err = GetWebShopLocations(x)
	if err == nil {
		rep = WebShopLocationsToProto(webShopLocations)
	}
	return rep, expires, err
}

func SetWebShopLocations(
	x cache.Context,
	rep map[b.LocationId]b.WebShopLocation,
) (
	err error,
) {
	return set(
		x,
		client.WriteWebShopLocations,
		keys.CacheKeyWebShopLocations,
		keys.TypeStrWebShopLocations,
		WEB_SHOP_LOCATIONS_EXPIRES_IN,
		rep,
		nil,
	)
}

func ProtoMergeSetWebShopLocations(
	x cache.Context,
	updates map[int64]*proto.CfgShopLocation,
) (
	err error,
) {
	if len(updates) == 0 {
		return nil
	}
	return protoMergeSetTerritories(
		x,
		updates,
		GetWebShopBundleKeys,
		GetWebShopLocations,
		ProtoMergeShopLocations,
		SetWebShopLocations,
	)
}

// // To Proto

func WebShopLocationsToProto(
	webShopLocations map[b.LocationId]b.WebShopLocation,
) (
	pbShopLocations map[int64]*proto.CfgShopLocation,
) {
	return newPBCfgShopLocations(webShopLocations)
}

func newPBCfgShopLocations(
	webShopLocations map[b.LocationId]b.WebShopLocation,
) (
	pbShopLocations map[int64]*proto.CfgShopLocation,
) {
	pbShopLocations = make(
		map[int64]*proto.CfgShopLocation,
		len(webShopLocations),
	)
	for locationId, webShopLocation := range webShopLocations {
		pbShopLocations[locationId] =
			newPBCfgShopLocation(webShopLocation)
	}
	return pbShopLocations
}

func newPBCfgShopLocation(
	webShopLocation b.WebShopLocation,
) (
	pbShopLocation *proto.CfgShopLocation,
) {
	return &proto.CfgShopLocation{
		BundleKey:   webShopLocation.BundleKey,
		TaxRate:     webShopLocation.TaxRate,
		BannedFlags: webShopLocation.BannedFlags,
	}
}

// // Merge

func ProtoMergeShopLocations(
	original map[b.LocationId]b.WebShopLocation,
	updates map[int64]*proto.CfgShopLocation,
	bundleKeys map[string]struct{},
) error {
	for locationId, pbShopLocation := range updates {
		if pbShopLocation == nil || pbShopLocation.BundleKey == "" {
			delete(original, locationId)
		} else if _, ok := bundleKeys[pbShopLocation.BundleKey]; !ok {
			return newPBtoWebShopLocationError(
				locationId,
				fmt.Sprintf(
					"type map key '%s' does not exist",
					pbShopLocation.BundleKey,
				),
			)
		} else {
			original[locationId] = pBtoWebShopLocation(
				pbShopLocation,
			)
		}
	}
	return nil
}

func newPBtoWebShopLocationError(
	locationId int64,
	errStr string,
) configerror.ErrInvalid {
	return configerror.ErrInvalid{
		Err: configerror.ErrShopLocationInvalid{
			Err: fmt.Errorf(
				"'%d': %s",
				locationId,
				errStr,
			),
		},
	}
}

func pBtoWebShopLocation(
	pbShopLocation *proto.CfgShopLocation,
) (
	webShopLocation b.WebShopLocation,
) {
	return b.WebShopLocation{
		BundleKey:   pbShopLocation.BundleKey,
		TaxRate:     pbShopLocation.TaxRate,
		BannedFlags: pbShopLocation.BannedFlags,
	}
}
