package proto

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
)

type PartialCfgBuybackSystemTypeMapsBuilderResponse struct {
	Systems         map[int32]*proto.CfgBuybackSystem
	SystemRegionMap map[int32]int32
}

type CfgGetBuybackSystemTypeMapsBuilderParams struct{}

type CfgGetBuybackSystemTypeMapsBuilderClient struct {
	webBTypeMapsBuilderReaderClient bucket.WebBuybackSystemTypeMapsBuilderReaderClient
}

func (gbsc CfgGetBuybackSystemTypeMapsBuilderClient) Fetch(
	ctx context.Context,
	params CfgGetBuybackSystemTypeMapsBuilderParams,
) (
	rep map[int32]*proto.CfgBuybackSystemTypeBundle,
	err error,
) {
	webBuilder, err := gbsc.fetchBuilder(ctx)
	if err != nil {
		return nil, err
	} else {
		return protoutil.NewPBCfgBTypeMapsBuilder(webBuilder), nil
	}
}

func (gbbc CfgGetBuybackSystemTypeMapsBuilderClient) fetchBuilder(
	ctx context.Context,
) (
	builder map[b.TypeId]b.WebBuybackSystemTypeBundle,
	err error,
) {
	if builderRep, err := gbbc.webBTypeMapsBuilderReaderClient.Fetch(
		ctx,
		bucket.WebBuybackSystemTypeMapsBuilderReaderParams{},
	); err != nil {
		return nil, err
	} else {
		return builderRep.Data(), nil
	}
}
