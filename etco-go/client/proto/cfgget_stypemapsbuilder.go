package proto

import (
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

type PartialCfgShopLocationTypeMapsBuilderResponse struct {
	Systems         map[int32]*proto.CfgShopLocation
	SystemRegionMap map[int32]int32
}

type CfgGetShopLocationTypeMapsBuilderParams struct{}

type CfgGetShopLocationTypeMapsBuilderClient struct{}

func NewCfgGetShopLocationTypeMapsBuilderClient() CfgGetShopLocationTypeMapsBuilderClient {
	return CfgGetShopLocationTypeMapsBuilderClient{}
}

func (gbsc CfgGetShopLocationTypeMapsBuilderClient) Fetch(
	x cache.Context,
	params CfgGetShopLocationTypeMapsBuilderParams,
) (
	rep map[int32]*proto.CfgShopLocationTypeBundle,
	err error,
) {
	webBuilder, err := gbsc.fetchBuilder(x)
	if err != nil {
		return nil, err
	} else {
		return protoutil.NewPBCfgSTypeMapsBuilder(webBuilder), nil
	}
}

func (gbbc CfgGetShopLocationTypeMapsBuilderClient) fetchBuilder(
	x cache.Context,
) (
	builder map[b.TypeId]b.WebShopLocationTypeBundle,
	err error,
) {
	builder, _, err = bucket.GetWebShopLocationTypeMapsBuilder(x)
	return builder, err
}
