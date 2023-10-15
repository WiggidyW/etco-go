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

type CfgMergeShopLocationsParams struct {
	Updates map[int64]*proto.CfgShopLocation
}

type CfgMergeShopLocationsClient struct {
	webShopLocationsReaderClient bucket.SC_WebShopLocationsReaderClient
	webShopLocationsWriterClient bucket.SAC_WebShopLocationsWriterClient
	webShopBundleKeysClient      bucket.SC_WebShopBundleKeysClient
}

func NewCfgMergeShopLocationsClient(
	webShopLocationsReaderClient bucket.SC_WebShopLocationsReaderClient,
	webShopLocationsWriterClient bucket.SAC_WebShopLocationsWriterClient,
	webShopBundleKeysClient bucket.SC_WebShopBundleKeysClient,
) CfgMergeShopLocationsClient {
	return CfgMergeShopLocationsClient{
		webShopLocationsReaderClient,
		webShopLocationsWriterClient,
		webShopBundleKeysClient,
	}
}

func (mslc CfgMergeShopLocationsClient) Fetch(
	ctx context.Context,
	params CfgMergeShopLocationsParams,
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
	go mslc.transceiveFetchBundleKeyHashSet(ctx, chanSendBundleKeyHashSet)

	// fetch the original systems
	systems, err := mslc.fetchSystems(ctx)
	if err != nil {
		return nil, err
	}

	// wait for the active bundle keys
	bundleKeyHashSet, err := chanRecvBundleKeyHashSet.Recv()
	if err != nil {
		return nil, err
	}

	// mutate the original systems with the updates
	if err = mergeShopLocations(
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
	if err = mslc.fetchWriteUpdated(ctx, systems); err != nil {
		return nil, err
	}

	return &CfgMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (mslc CfgMergeShopLocationsClient) fetchWriteUpdated(
	ctx context.Context,
	updated map[b.LocationId]b.WebShopLocation,
) error {
	_, err := mslc.webShopLocationsWriterClient.Fetch(
		ctx,
		bucket.WebShopLocationsWriterParams{
			WebShopLocations: updated,
		},
	)
	return err
}

func (mslc CfgMergeShopLocationsClient) fetchSystems(
	ctx context.Context,
) (
	systems map[b.LocationId]b.WebShopLocation,
	err error,
) {
	systemsRep, err := mslc.webShopLocationsReaderClient.Fetch(
		ctx,
		bucket.WebShopLocationsReaderParams{},
	)
	if err != nil {
		return nil, err
	} else {
		return systemsRep.Data(), nil
	}
}

func (mslc CfgMergeShopLocationsClient) transceiveFetchBundleKeyHashSet(
	ctx context.Context,
	chnSend chanresult.ChanSendResult[util.MapHashSet[string, struct{}]],
) error {
	bundleKeyHashSet, err := mslc.fetchBundleKeyHashSet(ctx)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(bundleKeyHashSet)
	}
}

func (mslc CfgMergeShopLocationsClient) fetchBundleKeyHashSet(
	ctx context.Context,
) (
	bundleKeyHashSet util.MapHashSet[string, struct{}],
	err error,
) {
	bundleKeys, err := mslc.webShopBundleKeysClient.Fetch(
		ctx,
		bucket.WebShopBundleKeysParams{},
	)
	if err != nil {
		return nil, err
	} else {
		return util.MapHashSet[string, struct{}](bundleKeys.Data()), nil
	}
}

func mergeShopLocations[HS util.HashSet[string]](
	original map[b.LocationId]b.WebShopLocation,
	updates map[int64]*proto.CfgShopLocation,
	bundleKeys HS,
) error {
	for locationId, pbShopLocation := range updates {
		if pbShopLocation == nil || pbShopLocation.BundleKey == "" {
			delete(original, locationId)
		} else if !bundleKeys.Has(pbShopLocation.BundleKey) {
			return newPBtoWebShopLocationError(
				locationId,
				fmt.Sprintf(
					"type map key '%s' does not exist",
					pbShopLocation.BundleKey,
				),
			)
		} else {
			original[locationId] = pBtoWebShopLocation(
				pbShopLocation,
			)
		}
	}
	return nil
}

func newPBtoWebShopLocationError(
	locationId int64,
	errStr string,
) configerror.ErrInvalid {
	return configerror.ErrInvalid{
		Err: configerror.ErrShopLocationInvalid{
			Err: fmt.Errorf(
				"'%d': %s",
				locationId,
				errStr,
			),
		},
	}
}

func pBtoWebShopLocation(
	pbShopLocation *proto.CfgShopLocation,
) (
	webShopLocation b.WebShopLocation,
) {
	return b.WebShopLocation{
		BundleKey:   pbShopLocation.BundleKey,
		BannedFlags: pbShopLocation.BannedFlags,
	}
}
