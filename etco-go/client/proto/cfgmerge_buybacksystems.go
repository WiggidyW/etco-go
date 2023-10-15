package proto

import (
	"context"
	"fmt"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/client/bucket"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/util"
)

type CfgMergeBuybackSystemsParams struct {
	Updates map[int32]*proto.CfgBuybackSystem
}

type CfgMergeBuybackSystemsClient struct {
	webBuybackSystemsReaderClient bucket.SC_WebBuybackSystemsReaderClient
	webBuybackSystemsWriterClient bucket.SAC_WebBuybackSystemsWriterClient
	webBuybackBundleKeysClient    bucket.SC_WebBuybackBundleKeysClient
}

func NewCfgMergeBuybackSystemsClient(
	webBuybackSystemsReaderClient bucket.SC_WebBuybackSystemsReaderClient,
	webBuybackSystemsWriterClient bucket.SAC_WebBuybackSystemsWriterClient,
	webBuybackBundleKeysClient bucket.SC_WebBuybackBundleKeysClient,
) CfgMergeBuybackSystemsClient {
	return CfgMergeBuybackSystemsClient{
		webBuybackSystemsReaderClient,
		webBuybackSystemsWriterClient,
		webBuybackBundleKeysClient,
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
	bundleKeys, err := mbsc.webBuybackBundleKeysClient.Fetch(
		ctx,
		bucket.WebBuybackBundleKeysParams{},
	)
	if err != nil {
		return nil, err
	} else {
		return util.MapHashSet[string, struct{}](bundleKeys.Data()), nil
	}
}

func mergeBuybackSystems[HS util.HashSet[string]](
	original map[b.SystemId]b.WebBuybackSystem,
	updates map[int32]*proto.CfgBuybackSystem,
	bundleKeys HS,
) error {
	for systemId, pbBuybackSystem := range updates {
		if pbBuybackSystem == nil || pbBuybackSystem.BundleKey == "" {
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
		TaxRate:   pbBuybackSystem.TaxRate,
		M3Fee:     pbBuybackSystem.M3Fee,
	}
}
