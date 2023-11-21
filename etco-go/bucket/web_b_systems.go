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
	WEB_BUYBACK_SYSTEMS_BUF_CAP    int           = 0
	WEB_BUYBACK_SYSTEMS_EXPIRES_IN time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebBuybackSystems = cache.RegisterType[map[b.SystemId]b.WebBuybackSystem]("webbuybacksystems", WEB_BUYBACK_SYSTEMS_BUF_CAP)
}

func GetWebBuybackSystems(
	x cache.Context,
) (
	rep map[b.SystemId]b.WebBuybackSystem,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebBuybackSystems,
		keys.CacheKeyWebBuybackSystems,
		keys.TypeStrWebBuybackSystems,
		WEB_BUYBACK_SYSTEMS_EXPIRES_IN,
		build.CAPACITY_WEB_BUYBACK_SYSTEMS,
	)
}

func ProtoGetWebBuybackSystems(
	x cache.Context,
) (
	rep map[int32]*proto.CfgBuybackSystem,
	expires time.Time,
	err error,
) {
	var webBuybackSystems map[b.SystemId]b.WebBuybackSystem
	webBuybackSystems, expires, err = GetWebBuybackSystems(x)
	if err == nil {
		rep = WebBuybackSystemsToProto(webBuybackSystems)
	}
	return rep, expires, err
}

func SetWebBuybackSystems(
	x cache.Context,
	rep map[b.SystemId]b.WebBuybackSystem,
) (
	err error,
) {
	return set(
		x,
		client.WriteWebBuybackSystems,
		keys.CacheKeyWebBuybackSystems,
		keys.TypeStrWebBuybackSystems,
		WEB_BUYBACK_SYSTEMS_EXPIRES_IN,
		rep,
		nil,
	)
}

func ProtoMergeSetWebBuybackSystems(
	x cache.Context,
	updates map[int32]*proto.CfgBuybackSystem,
) (
	err error,
) {
	if len(updates) == 0 {
		return nil
	}
	return protoMergeSetTerritories(
		x,
		updates,
		GetWebBuybackBundleKeys,
		GetWebBuybackSystems,
		ProtoMergeBuybackSystems,
		SetWebBuybackSystems,
	)
}

// // To Proto

func WebBuybackSystemsToProto(
	webBuybackSystems map[b.SystemId]b.WebBuybackSystem,
) (
	pbBuybackSystems map[int32]*proto.CfgBuybackSystem,
) {
	return newPBCfgBuybackSystems(webBuybackSystems)
}

func newPBCfgBuybackSystems(
	webBuybackSystems map[b.SystemId]b.WebBuybackSystem,
) (
	pbBuybackSystems map[int32]*proto.CfgBuybackSystem,
) {
	pbBuybackSystems = make(
		map[int32]*proto.CfgBuybackSystem,
		len(webBuybackSystems),
	)
	for systemId, webBuybackSystem := range webBuybackSystems {
		pbBuybackSystems[systemId] =
			newPBCfgBuybackSystem(webBuybackSystem)
	}
	return pbBuybackSystems
}

func newPBCfgBuybackSystem(
	webBuybackSystem b.WebBuybackSystem,
) (
	pbBuybackSystem *proto.CfgBuybackSystem,
) {
	return &proto.CfgBuybackSystem{
		BundleKey: webBuybackSystem.BundleKey,
		TaxRate:   webBuybackSystem.TaxRate,
		M3Fee:     webBuybackSystem.M3Fee,
	}
}

// // Merge

func ProtoMergeBuybackSystems(
	original map[b.SystemId]b.WebBuybackSystem,
	updates map[int32]*proto.CfgBuybackSystem,
	bundleKeys map[string]struct{},
) error {
	for systemId, pbBuybackSystem := range updates {
		if pbBuybackSystem == nil || pbBuybackSystem.BundleKey == "" {
			delete(original, systemId)
		} else if _, ok := bundleKeys[pbBuybackSystem.BundleKey]; !ok {
			return newPBtoWebBuybackSystemError(
				systemId,
				fmt.Sprintf(
					"type map key '%s' does not exist",
					pbBuybackSystem.BundleKey,
				),
			)
		} else {
			original[systemId] = pBtoWebBuybackSystem(
				pbBuybackSystem,
			)
		}
	}
	return nil
}

func newPBtoWebBuybackSystemError(
	systemId int32,
	errStr string,
) configerror.ErrInvalid {
	return configerror.ErrInvalid{
		Err: configerror.ErrBuybackSystemInvalid{
			Err: fmt.Errorf(
				"'%d': %s",
				systemId,
				errStr,
			),
		},
	}
}

func pBtoWebBuybackSystem(
	pbBuybackSystem *proto.CfgBuybackSystem,
) (
	webBuybackSystem b.WebBuybackSystem,
) {
	return b.WebBuybackSystem{
		BundleKey: pbBuybackSystem.BundleKey,
		TaxRate:   pbBuybackSystem.TaxRate,
		M3Fee:     pbBuybackSystem.M3Fee,
	}
}
