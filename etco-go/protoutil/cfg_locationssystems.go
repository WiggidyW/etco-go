package protoutil

import (
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/proto"
)

func NewPBCfgBuybackSystems(
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
			NewPBCfgBuybackSystem(webBuybackSystem)
	}
	return pbBuybackSystems
}

func NewPBCfgBuybackSystem(
	webBuybackSystem b.WebBuybackSystem,
) (
	pbBuybackSystem *proto.CfgBuybackSystem,
) {
	return &proto.CfgBuybackSystem{
		BundleKey: webBuybackSystem.BundleKey,
		M3Fee:     webBuybackSystem.M3Fee,
	}
}

func NewPBCfgShopLocations(
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
			NewPBCfgShopLocation(webShopLocation)
	}
	return pbShopLocations
}

func NewPBCfgShopLocation(
	webShopLocation b.WebShopLocation,
) (
	pbShopLocation *proto.CfgShopLocation,
) {
	return &proto.CfgShopLocation{
		BundleKey:   webShopLocation.BundleKey,
		BannedFlags: webShopLocation.BannedFlags,
	}
}
