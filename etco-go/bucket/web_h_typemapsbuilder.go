package bucket

import (
	"fmt"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/proto"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_H_TYPEMAPSBUILDER_BUF_CAP    int           = 0
	WEB_H_TYPEMAPSBUILDER_EXPIRES_IN time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebHaulRouteTypeMapsBuilder = cache.RegisterType[map[b.TypeId]b.WebHaulRouteTypeBundle]("webhaulroutetypemapsbuilder", WEB_H_TYPEMAPSBUILDER_BUF_CAP)
}

func GetWebHaulRouteTypeMapsBuilder(
	x cache.Context,
) (
	rep map[b.TypeId]b.WebHaulRouteTypeBundle,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebHaulRouteTypeMapsBuilder,
		keys.CacheKeyWebHaulRouteTypeMapsBuilder,
		keys.TypeStrWebHaulRouteTypeMapsBuilder,
		WEB_H_TYPEMAPSBUILDER_EXPIRES_IN,
		build.CAPACITY_WEB_HAUL_ROUTE_TYPE_MAPS_BUILDER,
	)
}

func GetWebHaulRouteActiveMarkets(
	x cache.Context,
) (
	rep map[string]struct{},
	expires time.Time,
	err error,
) {
	var webHaulRouteTypeMapsBuilder map[b.TypeId]b.WebHaulRouteTypeBundle
	webHaulRouteTypeMapsBuilder, expires, err =
		GetWebHaulRouteTypeMapsBuilder(x)
	if err == nil {
		rep = extractHaulBuilderActiveMarkets(
			webHaulRouteTypeMapsBuilder,
		)
	}
	return rep, expires, err
}

func ProtoGetWebHaulRouteTypeMapsBuilder(
	x cache.Context,
) (
	rep map[int32]*proto.CfgHaulRouteTypeBundle,
	expires time.Time,
	err error,
) {
	var webHaulRouteTypeMapsBuilder map[b.TypeId]b.WebHaulRouteTypeBundle
	webHaulRouteTypeMapsBuilder, expires, err =
		GetWebHaulRouteTypeMapsBuilder(x)
	if err == nil {
		rep = WebHaulRouteTypeMapsBuilderToProto(
			webHaulRouteTypeMapsBuilder,
		)
	}
	return rep, expires, err
}

func SetWebHaulRouteTypeMapsBuilder(
	x cache.Context,
	rep map[b.TypeId]b.WebHaulRouteTypeBundle,
) (
	err error,
) {
	lock := cacheprefetch.ServerLock(
		keys.CacheKeyWebHaulBundleKeys,
		keys.TypeStrWebHaulBundleKeys,
	)
	return set(
		x,
		client.WriteWebHaulRouteTypeMapsBuilder,
		keys.CacheKeyWebHaulRouteTypeMapsBuilder,
		keys.TypeStrWebHaulRouteTypeMapsBuilder,
		WEB_H_TYPEMAPSBUILDER_EXPIRES_IN,
		rep,
		&lock,
	)
}

func ProtoMergeSetWebHaulRouteTypeMapsBuilder(
	x cache.Context,
	updates map[int32]*proto.CfgHaulRouteTypeBundle,
) (
	err error,
) {
	if len(updates) == 0 {
		return nil
	}
	return protoMergeSetTypeMapsBuilder(
		x,
		updates,
		GetWebHaulRouteTypeMapsBuilder,
		ProtoMergeHaulRouteTypeMapsBuilder,
		SetWebHaulRouteTypeMapsBuilder,
	)
}

// // Active Markets

func extractHaulBuilderActiveMarkets(
	sbuilder map[b.TypeId]b.WebHaulRouteTypeBundle,
) map[string]struct{} {
	activeMarkets := make(map[string]struct{})
	for _, bundle := range sbuilder {
		for _, hTypePricing := range bundle {
			activeMarkets[hTypePricing.MarketName] = struct{}{}
		}
	}
	return activeMarkets
}

// // To Proto

func WebHaulRouteTypeMapsBuilderToProto(
	webBTypeMapsBuilder map[b.TypeId]b.WebHaulRouteTypeBundle,
) (
	pbBTypeMapsBuilder map[int32]*proto.CfgHaulRouteTypeBundle,
) {
	return newPBCfgHTypeMapsBuilder(webBTypeMapsBuilder)
}

func newPBCfgHTypeMapsBuilder(
	webHTypeMapsBuilder map[b.TypeId]b.WebHaulRouteTypeBundle,
) (
	pbHTypeMapsBuilder map[int32]*proto.CfgHaulRouteTypeBundle,
) {
	pbHTypeMapsBuilder = make(
		map[int32]*proto.CfgHaulRouteTypeBundle,
		len(webHTypeMapsBuilder),
	)
	for typeId, webTypeBundle := range webHTypeMapsBuilder {
		pbHTypeMapsBuilder[typeId] =
			newPBCfgHaulRouteTypeBundle(webTypeBundle)
	}
	return pbHTypeMapsBuilder
}

func newPBCfgHaulRouteTypeBundle(
	webHaulRouteTypeBundle b.WebHaulRouteTypeBundle,
) (
	pbHaulRouteTypeBundle *proto.CfgHaulRouteTypeBundle,
) {
	pbHaulRouteTypeBundle = &proto.CfgHaulRouteTypeBundle{
		Inner: make(
			map[string]*proto.CfgHaulTypePricing,
			len(webHaulRouteTypeBundle),
		),
	}
	for bundleKey, webHTypePricing := range webHaulRouteTypeBundle {
		pbHaulRouteTypeBundle.Inner[bundleKey] =
			&proto.CfgHaulTypePricing{
				Inner: newPBCfgTypePricing(webHTypePricing),
			}
	}
	return pbHaulRouteTypeBundle
}

// // Merge

func ProtoMergeHaulRouteTypeMapsBuilder[HSV any](
	original map[b.TypeId]b.WebHaulRouteTypeBundle,
	updates map[int32]*proto.CfgHaulRouteTypeBundle,
	markets map[string]HSV,
) error {
	for typeId, pbHaulRouteTypeBundle := range updates {
		if pbHaulRouteTypeBundle == nil ||
			pbHaulRouteTypeBundle.Inner == nil ||
			len(pbHaulRouteTypeBundle.Inner) == 0 {
			delete(original, typeId)
		} else {
			haulRouteTypeBundle, ok := original[typeId]
			if !ok {
				haulRouteTypeBundle = make(
					b.WebHaulRouteTypeBundle,
					len(pbHaulRouteTypeBundle.Inner),
				)
				original[typeId] = haulRouteTypeBundle
			}
			if err := mergeHaulRouteTypeBundle(
				typeId,
				haulRouteTypeBundle,
				pbHaulRouteTypeBundle.Inner,
				markets,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func mergeHaulRouteTypeBundle[HSV any](
	typeId b.TypeId,
	original b.WebHaulRouteTypeBundle,
	updates map[string]*proto.CfgHaulTypePricing,
	markets map[string]HSV,
) error {
	for bundleKey, pbHaulTypePricing := range updates {
		if pbHaulTypePricing == nil || pbHaulTypePricing.Inner == nil {
			delete(original, bundleKey)
		} else {
			haulTypePricing, err := pBToWebHaulTypePricing(
				typeId,
				bundleKey,
				pbHaulTypePricing.Inner,
				markets,
			)
			if err != nil {
				return err
			} else {
				original[bundleKey] = *haulTypePricing
			}
		}
	}
	return nil
}

func newPBtoWebHaulTypePricingError(
	typeId int32,
	typeMapKey string,
	errStr string,
) configerror.ErrInvalid {
	return configerror.ErrInvalid{
		Err: configerror.ErrHaulTypeInvalid{
			Err: fmt.Errorf(
				"'%d - %s': %s",
				typeId,
				typeMapKey,
				errStr,
			),
		},
	}
}

func pBToWebHaulTypePricing[HSV any](
	typeId b.TypeId,
	bundleKey b.BundleKey,
	pbHaulTypePricing *proto.CfgTypePricing,
	markets map[string]HSV,
) (
	webHaulTypePricing *b.WebHaulRouteTypePricing,
	err error,
) {
	webHaulTypePricing, err = PBtoWebTypePricing(
		pbHaulTypePricing,
		markets,
	)

	if err != nil {
		return nil, newPBtoWebHaulTypePricingError(
			typeId,
			bundleKey,
			err.Error(),
		)
	}

	return webHaulTypePricing, nil
}
