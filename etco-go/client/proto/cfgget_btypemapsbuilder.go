package proto

import (
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

type PartialCfgBuybackSystemTypeMapsBuilderResponse struct {
	Systems         map[int32]*proto.CfgBuybackSystem
	SystemRegionMap map[int32]int32
}

type CfgGetBuybackSystemTypeMapsBuilderParams struct{}

type CfgGetBuybackSystemTypeMapsBuilderClient struct{}

func NewCfgGetBuybackSystemTypeMapsBuilderClient() CfgGetBuybackSystemTypeMapsBuilderClient {
	return CfgGetBuybackSystemTypeMapsBuilderClient{}
}

func (gbsc CfgGetBuybackSystemTypeMapsBuilderClient) Fetch(
	x cache.Context,
	params CfgGetBuybackSystemTypeMapsBuilderParams,
) (
	rep map[int32]*proto.CfgBuybackSystemTypeBundle,
	err error,
) {
	webBuilder, err := gbsc.fetchBuilder(x)
	if err != nil {
		return nil, err
	} else {
		return protoutil.NewPBCfgBTypeMapsBuilder(webBuilder), nil
	}
}

func (gbbc CfgGetBuybackSystemTypeMapsBuilderClient) fetchBuilder(
	x cache.Context,
) (
	builder map[b.TypeId]b.WebBuybackSystemTypeBundle,
	err error,
) {
	builder, _, err = bucket.GetWebBuybackSystemTypeMapsBuilder(x)
	return builder, err
}
