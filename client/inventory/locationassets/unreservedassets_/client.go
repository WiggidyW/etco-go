package unreservedassets_

import (
	"context"
	"time"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/cache"
	sc "github.com/WiggidyW/etco-go/client/caching/strong/caching"
	massetscorporation "github.com/WiggidyW/etco-go/client/esi/model/assetscorporation"
	"github.com/WiggidyW/etco-go/client/inventory/locationassets/unreservedassets_/allassets_"
	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

const (
	UNRESERVED_SHOP_ASSETS_MIN_EXPIRES   time.Duration = 0
	UNRESERVED_SHOP_ASSETS_SLOCK_TTL     time.Duration = 1 * time.Minute
	UNRESERVED_SHOP_ASSETS_SLOCK_MAXWAIT time.Duration = 1 * time.Minute
)

type SC_UnreservedShopAssetsClient = sc.StrongCachingClient[
	UnreservedShopAssetsParams,
	map[int64]map[int32]int64,
	cache.ExpirableData[map[int64]map[int32]int64],
	UnreservedShopAssetsClient,
]

func NewSC_UnreservedShopAssetsClient(
	modelacClient massetscorporation.AssetsCorporationClient,
	appraisalClient rdbc.WC_ReadShopAppraisalClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) SC_UnreservedShopAssetsClient {
	return sc.NewStrongCachingClient(
		NewUnreservedShopAssetsClient(
			modelacClient,
			appraisalClient,
			cCache,
			sCache,
		),
		UNRESERVED_SHOP_ASSETS_MIN_EXPIRES,
		sCache,
		UNRESERVED_SHOP_ASSETS_SLOCK_TTL,
		UNRESERVED_SHOP_ASSETS_SLOCK_MAXWAIT,
	)
}

type UnreservedShopAssetsClient struct {
	allClient       allassets_.WC_AllShopAssetsClient
	appraisalClient rdbc.WC_ReadShopAppraisalClient
}

func NewUnreservedShopAssetsClient(
	modelacClient massetscorporation.AssetsCorporationClient,
	appraisalClient rdbc.WC_ReadShopAppraisalClient,
	cCache cache.SharedClientCache,
	sCache cache.SharedServerCache,
) UnreservedShopAssetsClient {
	return UnreservedShopAssetsClient{
		allClient: allassets_.NewWC_AllShopAssetsClient(
			modelacClient,
			cCache,
			sCache,
		),
		appraisalClient: appraisalClient,
	}
}

func (usac UnreservedShopAssetsClient) Fetch(
	ctx context.Context,
	params UnreservedShopAssetsParams,
) (*cache.ExpirableData[map[int64]map[int32]int64], error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// send out fetches for the appraisals keyed by the shop queue codes
	chnSend, chnRecv := chanresult.
		NewChanResult[*rdb.ShopAppraisal](ctx, 0, 0).Split()
	for _, code := range params.ShopQueue {
		go usac.fetchAppraisal(ctx, code, chnSend)
	}

	// fetch the assets
	assetsRep, err := usac.allClient.Fetch(
		ctx,
		allassets_.AllShopAssetsParams{},
	)
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
	chnSend chanresult.ChanSendResult[*rdb.ShopAppraisal],
) {
	if rep, err := usac.appraisalClient.Fetch(
		ctx,
		rdbc.ReadAppraisalParams{AppraisalCode: code},
	); err != nil {
		chnSend.SendErr(err)
	} else {
		chnSend.SendOk(rep.Data())
	}
}

func filterReserved(
	assets map[int64]map[int32]*int64,
	appraisal rdb.ShopAppraisal,
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
