package proto

import (
	"context"
	"fmt"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/util"
)

type CfgMergeBuybackSystemsParams struct {
	Updates map[int32]*proto.CfgBuybackSystem
}

type CfgMergeBuybackSystemsClient struct {
	webBuybackSystemsReaderClient   bucket.SC_WebBuybackSystemsReaderClient
	webBuybackSystemsWriterClient   bucket.SAC_WebBuybackSystemsWriterClient
	webBTypeMapsBuilderReaderClient bucket.SC_WebBuybackSystemTypeMapsBuilderReaderClient
	webSTypeMapsBuilderReaderClient bucket.SC_WebShopLocationTypeMapsBuilderReaderClient
}

func NewCfgMergeBuybackSystemsClient(
	webBuybackSystemsReaderClient bucket.SC_WebBuybackSystemsReaderClient,
	webBuybackSystemsWriterClient bucket.SAC_WebBuybackSystemsWriterClient,
	webBTypeMapsBuilderReaderClient bucket.SC_WebBuybackSystemTypeMapsBuilderReaderClient,
	webSTypeMapsBuilderReaderClient bucket.SC_WebShopLocationTypeMapsBuilderReaderClient,
) CfgMergeBuybackSystemsClient {
	return CfgMergeBuybackSystemsClient{
		webBuybackSystemsReaderClient,
		webBuybackSystemsWriterClient,
		webBTypeMapsBuilderReaderClient,
		webSTypeMapsBuilderReaderClient,
	}
}

func (mbsc CfgMergeBuybackSystemsClient) Fetch(
	ctx context.Context,
	params CfgMergeBuybackSystemsParams,
) (
	rep *CfgMergeResponse,
	err error,
) {
	// if there are no updates, return now
	if params.Updates == nil || len(params.Updates) == 0 {
		return &CfgMergeResponse{
			// Modified: false,
			// MergeError: nil,
		}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the active bundle keys for both buyback and shop in a goroutine
	chanSendBundleKeyHashSet, chanRecvBundleKeyHashSet := chanresult.
		NewChanResult[util.MapHashSet[string, struct{}]](
		ctx, 0, 0,
	).Split()
	go mbsc.transceiveFetchBundleKeyHashSet(ctx, chanSendBundleKeyHashSet)

	// fetch the original systems
	systems, err := mbsc.fetchSystems(ctx)
	if err != nil {
		return nil, err
	}

	// wait for the active bundle keys
	bundleKeyHashSet, err := chanRecvBundleKeyHashSet.Recv()
	if err != nil {
		return nil, err
	}

	// mutate the original systems with the updates
	if err = mergeBuybackSystems(
		systems,
		params.Updates,
		bundleKeyHashSet,
	); err != nil {
		return &CfgMergeResponse{
			// Modified: false,
			MergeError: err,
		}, nil
	}

	// write the mutated systems
	if err = mbsc.fetchWriteUpdated(ctx, systems); err != nil {
		return nil, err
	}

	return &CfgMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (mbsc CfgMergeBuybackSystemsClient) fetchWriteUpdated(
	ctx context.Context,
	updated map[b.SystemId]b.WebBuybackSystem,
) error {
	_, err := mbsc.webBuybackSystemsWriterClient.Fetch(
		ctx,
		bucket.WebBuybackSystemsWriterParams{
			WebBuybackSystems: updated,
		},
	)
	return err
}

func (mbsc CfgMergeBuybackSystemsClient) fetchSystems(
	ctx context.Context,
) (
	systems map[b.SystemId]b.WebBuybackSystem,
	err error,
) {
	systemsRep, err := mbsc.webBuybackSystemsReaderClient.Fetch(
		ctx,
		bucket.WebBuybackSystemsReaderParams{},
	)
	if err != nil {
		return nil, err
	} else {
		return systemsRep.Data(), nil
	}
}

func (mbsc CfgMergeBuybackSystemsClient) transceiveFetchBundleKeyHashSet(
	ctx context.Context,
	chnSend chanresult.ChanSendResult[util.MapHashSet[string, struct{}]],
) error {
	bundleKeyHashSet, err := mbsc.fetchBundleKeyHashSet(ctx)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(bundleKeyHashSet)
	}
}

func (mbsc CfgMergeBuybackSystemsClient) fetchBundleKeyHashSet(
	ctx context.Context,
) (
	bundleKeyHashSet util.MapHashSet[string, struct{}],
	err error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSendBundleKeys, chnRecvBundleKeys := chanresult.
		NewChanResult[map[string]struct{}](ctx, 1, 0).Split()

	// fetch the bigger one locally, and spawn a goroutine for the smaller one
	// we already know which is bigger from build constants
	var bigBundleKeys map[string]struct{}
	if build.CAPACITY_CORE_BUYBACK_SYSTEM_TYPE_MAPS >
		build.CAPACITY_CORE_SHOP_LOCATION_TYPE_MAPS {
		go mbsc.transceiveFetchSBuilderBundleKeys(
			ctx,
			chnSendBundleKeys,
		)
		bigBundleKeys, err = mbsc.fetchBBuilderBundleKeys(ctx)
	} else {
		go mbsc.transceiveFetchBBuilderBundleKeys(
			ctx,
			chnSendBundleKeys,
		)
		bigBundleKeys, err = mbsc.fetchSBuilderBundleKeys(ctx)
	}
	if err != nil {
		return nil, err
	}

	smallBundleKeys, err := chnRecvBundleKeys.Recv()
	if err != nil {
		return nil, err
	}

	// insert the small bundle keys into the big bundle keys
	for bundleKey := range smallBundleKeys {
		bigBundleKeys[bundleKey] = struct{}{}
	}

	return util.MapHashSet[string, struct{}](bigBundleKeys), nil
}

func (mbsc CfgMergeBuybackSystemsClient) transceiveFetchBBuilderBundleKeys(
	ctx context.Context,
	chnSend chanresult.ChanSendResult[map[string]struct{}],
) error {
	bundleKeys, err := mbsc.fetchBBuilderBundleKeys(ctx)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(bundleKeys)
	}
}

func (mbsc CfgMergeBuybackSystemsClient) fetchBBuilderBundleKeys(
	ctx context.Context,
) (
	bundleKeys map[string]struct{},
	err error,
) {
	bBuilderRep, err := mbsc.webBTypeMapsBuilderReaderClient.Fetch(
		ctx,
		bucket.WebBuybackSystemTypeMapsBuilderReaderParams{},
	)
	if err != nil {
		return nil, err
	}
	return extractBuilderBundleKeys(bBuilderRep.Data()), nil
}

func (mbsc CfgMergeBuybackSystemsClient) transceiveFetchSBuilderBundleKeys(
	ctx context.Context,
	chnSend chanresult.ChanSendResult[map[string]struct{}],
) error {
	bundleKeys, err := mbsc.fetchSBuilderBundleKeys(ctx)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(bundleKeys)
	}
}

func (mbsc CfgMergeBuybackSystemsClient) fetchSBuilderBundleKeys(
	ctx context.Context,
) (
	bundleKeys map[string]struct{},
	err error,
) {
	sBuilderRep, err := mbsc.webSTypeMapsBuilderReaderClient.Fetch(
		ctx,
		bucket.WebShopLocationTypeMapsBuilderReaderParams{},
	)
	if err != nil {
		return nil, err
	}
	return extractBuilderBundleKeys(sBuilderRep.Data()), nil
}

func mergeBuybackSystems[HS util.HashSet[string]](
	original map[b.SystemId]b.WebBuybackSystem,
	updates map[int32]*proto.CfgBuybackSystem,
	bundleKeys HS,
) error {
	for systemId, pbBuybackSystem := range updates {
		if pbBuybackSystem == nil {
			delete(original, systemId)
		} else if !bundleKeys.Has(pbBuybackSystem.BundleKey) {
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
		M3Fee:     pbBuybackSystem.M3Fee,
	}
}
