package protoutil

import (
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
)

func NewPBBuybackSystems(
	infoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker],
) (
	buybackSystems []*proto.BuybackSystem,
) {
	UNSAFE_BuybackSystems := staticdb.UnsafeGetCoreBuybackSystems()
	buybackSystems = make(
		[]*proto.BuybackSystem,
		0,
		len(UNSAFE_BuybackSystems),
	)

	for systemId := range UNSAFE_BuybackSystems {
		regionId := MaybeAddSystem(infoSession, systemId)
		buybackSystems = append(buybackSystems, &proto.BuybackSystem{
			SystemId: systemId,
			RegionId: regionId,
		})
	}

	return buybackSystems
}
