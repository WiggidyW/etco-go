package unreserved

import (
	"context"

	"github.com/WiggidyW/weve-esi/cache"
	sc "github.com/WiggidyW/weve-esi/client/caching/strong/caching"
	"github.com/WiggidyW/weve-esi/client/inventory/internal/all"
	"github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	"github.com/WiggidyW/weve-esi/client/remotedb/appraisal/readshop"
	"github.com/WiggidyW/weve-esi/util"
)

type SC_UnreservedShopAssetsClient = sc.StrongCachingClient[
	UnreservedShopAssetsParams,
	map[int64]map[int32]int64,
	cache.ExpirableData[map[int64]map[int32]int64],
	UnreservedShopAssetsClient,
]

type UnreservedShopAssetsClient struct {
	allClient       all.WC_AllShopAssetsClient
	appraisalClient readshop.WC_ReadShopAppraisalClient
}

func (usac UnreservedShopAssetsClient) Fetch(
	ctx context.Context,
	params UnreservedShopAssetsParams,
) (*cache.ExpirableData[map[int64]map[int32]int64], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// send out fetches for the appraisals keyed by the shop queue codes
	chnSend, chnRecv := util.NewChanResult[*appraisal.DBShopAppraisal](ctx).
		Split()
	for _, code := range params.ShopQueue {
		go usac.fetchAppraisal(ctx, code, chnSend)
	}

	// fetch the assets
	assetsRep, err := usac.allClient.Fetch(ctx, all.AllShopAssetsParams{})
	if err != nil {
		return nil, err
	}
	shopAssets := assetsRep.Data()

	// // return early if there are no assets or no shop queue codes
	// if len(shopAssets) == 0 || len(params.ShopQueue) == 0 {
	// 	return cache.NewExpirableDataPtr(
	// 		shopAssets,
	// 		assetsRep.Expires(),
	// 	), nil
	// }

	// filter out the reserved appraisal items from the shop queue
	for i := 0; i < len(params.ShopQueue); i++ {
		appraisal, err := chnRecv.Recv()
		if err != nil {
			return nil, err
		} else if appraisal == nil {
			continue
		}

		filterReserved(shopAssets, *appraisal)
	}

	return cache.NewExpirableDataPtr(
		ptrMapToValMap(shopAssets),
		assetsRep.Expires(),
	), nil
}

func (usac UnreservedShopAssetsClient) fetchAppraisal(
	ctx context.Context,
	code string,
	chnSend util.ChanSendResult[*appraisal.DBShopAppraisal],
) {
	if rep, err := usac.appraisalClient.Fetch(
		ctx,
		appraisal.ReadAppraisalParams{AppraisalCode: code},
	); err != nil {
		chnSend.SendErr(err)
	} else {
		chnSend.SendOk(rep.Data())
	}
}

func filterReserved(
	assets map[int64]map[int32]*int64,
	appraisal appraisal.DBShopAppraisal,
) {
	locationAssets, ok := assets[appraisal.LocationId]
	if !ok {
		// return early if the location has no assets
		return
	}

	for _, item := range appraisal.Items {
		if assetQuantity, ok := locationAssets[item.TypeId]; ok {
			// subtract the quantity of the reserved item from the asset
			*assetQuantity -= item.Quantity
			// delete the asset if it's quantity is 0 or less
			if *assetQuantity <= 0 {
				delete(locationAssets, item.TypeId)
			}
		}
	}
}

func ptrMapToValMap(
	ptrMap map[int64]map[int32]*int64,
) map[int64]map[int32]int64 {
	valMap := make(map[int64]map[int32]int64, len(ptrMap))
	for locationId, locationAssets := range ptrMap {
		valMap[locationId] = make(
			map[int32]int64,
			len(locationAssets),
		)
		for typeId, asset := range locationAssets {
			valMap[locationId][typeId] = *asset
		}
	}
	return valMap
}
