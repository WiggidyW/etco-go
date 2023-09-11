package protoutil

import (
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
)

func NewPBSystems[V any](
	infoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker],
	rSystems map[int32]V,
) (
	systems []*proto.System,
) {
	systems = make(
		[]*proto.System,
		0,
		len(rSystems),
	)

	for systemId := range rSystems {
		regionId := MaybeAddSystem(infoSession, systemId)
		systems = append(systems, &proto.System{
			SystemId: systemId,
			RegionId: regionId,
		})
	}

	return systems
}

func NewPBSDESystems(
	infoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker],
) (
	systems []*proto.System,
) {
	UNSAFE_Systems := staticdb.UnsafeGetSDESystems()
	return NewPBSystems(infoSession, UNSAFE_Systems)
}

func NewPBBuybackSystems(
	infoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker],
) (
	buybackSystems []*proto.System,
) {
	UNSAFE_BuybackSystems := staticdb.UnsafeGetCoreBuybackSystems()
	return NewPBSystems(infoSession, UNSAFE_BuybackSystems)
}
