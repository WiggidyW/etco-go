package proto

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

type PartialCfgShopLocationTypeMapsBuilderResponse struct {
	Systems         map[int32]*proto.CfgShopLocation
	SystemRegionMap map[int32]int32
}

type CfgGetShopLocationTypeMapsBuilderParams struct{}

type CfgGetShopLocationTypeMapsBuilderClient struct {
	webSTypeMapsBuilderReaderClient bucket.SC_WebShopLocationTypeMapsBuilderReaderClient
}

func NewCfgGetShopLocationTypeMapsBuilderClient(
	webSTypeMapsBuilderReaderClient bucket.SC_WebShopLocationTypeMapsBuilderReaderClient,
) CfgGetShopLocationTypeMapsBuilderClient {
	return CfgGetShopLocationTypeMapsBuilderClient{
		webSTypeMapsBuilderReaderClient,
	}
}

func (gbsc CfgGetShopLocationTypeMapsBuilderClient) Fetch(
	ctx context.Context,
	params CfgGetShopLocationTypeMapsBuilderParams,
) (
	rep map[int32]*proto.CfgShopLocationTypeBundle,
	err error,
) {
	webBuilder, err := gbsc.fetchBuilder(ctx)
	if err != nil {
		return nil, err
	} else {
		return protoutil.NewPBCfgSTypeMapsBuilder(webBuilder), nil
	}
}

func (gbbc CfgGetShopLocationTypeMapsBuilderClient) fetchBuilder(
	ctx context.Context,
) (
	builder map[b.TypeId]b.WebShopLocationTypeBundle,
	err error,
) {
	if builderRep, err := gbbc.webSTypeMapsBuilderReaderClient.Fetch(
		ctx,
		bucket.WebShopLocationTypeMapsBuilderReaderParams{},
	); err != nil {
		return nil, err
	} else {
		return builderRep.Data(), nil
	}
}
