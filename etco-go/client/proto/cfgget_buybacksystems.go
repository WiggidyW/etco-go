package proto

import (
	"context"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

type PartialCfgBuybackSystemsResponse struct {
	Systems         map[int32]*proto.CfgBuybackSystem
	SystemRegionMap map[int32]int32
}

type CfgGetBuybackSystemsParams struct {
	LocationInfoSession *staticdb.LocationInfoSession[*staticdb.LocalLocationNamerTracker]
}

type CfgGetBuybackSystemsClient struct {
	webBuybackSystemsReaderClient bucket.SC_WebBuybackSystemsReaderClient
}

func NewCfgGetBuybackSystemsClient(
	webBuybackSystemsReaderClient bucket.SC_WebBuybackSystemsReaderClient,
) CfgGetBuybackSystemsClient {
	return CfgGetBuybackSystemsClient{webBuybackSystemsReaderClient}
}

func (gbsc CfgGetBuybackSystemsClient) Fetch(
	ctx context.Context,
	params CfgGetBuybackSystemsParams,
) (
	rep *PartialCfgBuybackSystemsResponse,
	err error,
) {
	// fetch web buyback systems
	webBuybackSystems, err := gbsc.fetchWebBuybackSystems(ctx)
	if err != nil {
		return nil, err
	}

	// if we don't need location info, convert it to PB and return now
	if params.LocationInfoSession == nil {
		return &PartialCfgBuybackSystemsResponse{
			Systems: protoutil.NewPBCfgBuybackSystems(
				webBuybackSystems,
			),
			// LocationInfoMap: nil,
		}, nil
	}

	rep = &PartialCfgBuybackSystemsResponse{
		Systems: make(
			map[int32]*proto.CfgBuybackSystem,
			len(webBuybackSystems),
		),
		SystemRegionMap: make(
			map[int32]int32,
			len(webBuybackSystems),
		),
	}

	for systemId, webSystem := range webBuybackSystems {
		rep.Systems[systemId] = protoutil.NewPBCfgBuybackSystem(
			webSystem,
		)
		rep.SystemRegionMap[systemId] = protoutil.MaybeAddSystem(
			params.LocationInfoSession,
			systemId,
		)
	}

	return rep, nil
}

func (gbsc CfgGetBuybackSystemsClient) fetchWebBuybackSystems(
	ctx context.Context,
) (
	buybackSystems map[b.SystemId]b.WebBuybackSystem,
	err error,
) {
	buybackSystemsRep, err := gbsc.webBuybackSystemsReaderClient.Fetch(
		ctx,
		bucket.WebBuybackSystemsReaderParams{},
	)
	if err != nil {
		return nil, err
	} else {
		return buybackSystemsRep.Data(), nil
	}
}
