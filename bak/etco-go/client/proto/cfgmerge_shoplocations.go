package proto

import (
	"fmt"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/util"
)

type CfgMergeShopLocationsParams struct {
	Updates map[int64]*proto.CfgShopLocation
}

type CfgMergeShopLocationsClient struct{}

func NewCfgMergeShopLocationsClient() CfgMergeShopLocationsClient {
	return CfgMergeShopLocationsClient{}
}

func (mslc CfgMergeShopLocationsClient) Fetch(
	x cache.Context,
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

	x, cancel := x.WithCancel()
	defer cancel()

	// fetch the active bundle keys for both buyback and shop in a goroutine
	chanSendBundleKeyHashSet, chanRecvBundleKeyHashSet := chanresult.
		NewChanResult[util.MapHashSet[string, struct{}]](
		x.Ctx(), 0, 0,
	).Split()
	go mslc.transceiveFetchBundleKeyHashSet(x, chanSendBundleKeyHashSet)

	// fetch the original locations
	locations, err := mslc.fetchLocations(x)
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
		locations,
		params.Updates,
		bundleKeyHashSet,
	); err != nil {
		return &CfgMergeResponse{
			// Modified: false,
			MergeError: err,
		}, nil
	}

	// write the mutated locations
	if err = mslc.fetchWriteUpdated(x, locations); err != nil {
		return nil, err
	}

	return &CfgMergeResponse{
		Modified: true,
		// MergeError: nil,
	}, nil
}

func (mslc CfgMergeShopLocationsClient) fetchWriteUpdated(
	x cache.Context,
	updated map[b.LocationId]b.WebShopLocation,
) error {
	return bucket.SetWebShopLocations(x, updated)
}

func (mslc CfgMergeShopLocationsClient) fetchLocations(
	x cache.Context,
) (
	locations map[b.LocationId]b.WebShopLocation,
	err error,
) {
	locations, _, err = bucket.GetWebShopLocations(x)
	return locations, err
}

func (mslc CfgMergeShopLocationsClient) transceiveFetchBundleKeyHashSet(
	x cache.Context,
	chnSend chanresult.ChanSendResult[util.MapHashSet[string, struct{}]],
) error {
	bundleKeyHashSet, err := mslc.fetchBundleKeyHashSet(x)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(bundleKeyHashSet)
	}
}

func (mslc CfgMergeShopLocationsClient) fetchBundleKeyHashSet(
	x cache.Context,
) (
	bundleKeyHashSet util.MapHashSet[string, struct{}],
	err error,
) {
	bundleKeys, _, err := bucket.GetWebShopBundleKeys(x)
	if err != nil {
		return nil, err
	} else {
		return util.MapHashSet[string, struct{}](bundleKeys), nil
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
		TaxRate:     pbShopLocation.TaxRate,
		BannedFlags: pbShopLocation.BannedFlags,
	}
}
