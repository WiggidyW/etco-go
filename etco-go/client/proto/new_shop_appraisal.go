package proto

import (
	"github.com/WiggidyW/etco-go/appraisal"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/items"
	"github.com/WiggidyW/etco-go/proto"
	pu "github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PBNewShopAppraisalParams[IM staticdb.IndexMap] struct {
	TypeNamingSession *staticdb.TypeNamingSession[IM]
	Items             []items.BasicItem
	LocationId        int64
	CharacterId       int32
	IncludeCode       bool
}

type PBNewShopAppraisalClient[IM staticdb.IndexMap] struct{}

func NewPBNewShopAppraisalClient[IM staticdb.IndexMap]() PBNewShopAppraisalClient[IM] {
	return PBNewShopAppraisalClient[IM]{}
}

func (nbac PBNewShopAppraisalClient[IM]) Fetch(
	x cache.Context,
	params PBNewShopAppraisalParams[IM],
) (*proto.ShopAppraisal, error) {
	rAppraisal, _, err := appraisal.CreateShopAppraisal(
		x,
		params.Items,
		&params.CharacterId,
		params.LocationId,
		params.IncludeCode,
	)
	if err != nil {
		return nil, err
	} else {
		return pu.NewPBShopAppraisal(
			rAppraisal,
			params.TypeNamingSession,
		), nil
	}
}
