package shopitems

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authing"
	"github.com/WiggidyW/weve-esi/client/smartclient/market"
	"github.com/WiggidyW/weve-esi/client/smartclient/shop/assets"
	"github.com/WiggidyW/weve-esi/staticdb"
	"github.com/WiggidyW/weve-esi/util"
)

type AuthingShopItemsClient = authing.AuthingClient[
	ShopItemsParams,
	ShopItemsResponse,
	ShopItemsClient,
]

type ShopItemsResponse struct {
	Items      []ShopItem
	MrktGroups []string
	Groups     []string
	Categories []string
}

type ShopItemsParams struct {
	UserRefreshToken  string // user's native refresh token
	CorporationId     int32
	CorpRefreshToken  string // corp's web refresh token
	LocationId        int64
	IncludeName       bool
	IncludeMrktGroups bool
	IncludeGroup      bool
	IncludeCategory   bool
}

func (sip ShopItemsParams) AuthRefreshToken() string {
	return sip.UserRefreshToken
}

type ShopItemsClient struct {
	assetClient assets.CachingShopLocationAssetsClient
	priceClient market.ShopClient
}

func (sic ShopItemsClient) Fetch(
	ctx context.Context,
	params ShopItemsParams,
) (*ShopItemsResponse, error) {
	// fetch the assets
	assetsRep, err := sic.assetClient.Fetch(
		ctx,
		assets.ShopLocationAssetsParams{
			CorporationId: params.CorporationId,
			RefreshToken:  params.CorpRefreshToken,
			LocationId:    params.LocationId,
		},
	)
	if err != nil {
		return nil, err
	}
	assets := assetsRep.Data()

	// return now if there are no assets at the provided location
	if len(assets) == 0 {
		return &ShopItemsResponse{}, nil
	}

	// // fetch the prices and names for the assets
	// start a naming session (for naming the items and deduping strings)
	namingSession := staticdb.NewNamingSession(
		params.IncludeName,
		params.IncludeMrktGroups,
		params.IncludeGroup,
		params.IncludeCategory,
	)
	// send out the fetches
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := util.NewChanResult[*ShopItem](ctx).Split()
	for _, asset := range assets {
		go sic.fetchItem(
			ctx,
			asset,
			params.LocationId,
			namingSession,
			chnSend,
		)
	}

	// collect the items
	items := make([]ShopItem, 0, len(assets))
	for i := 0; i < len(assets); i++ {
		if item, err := chnRecv.Recv(); err != nil {
			return nil, err
		} else if item != nil {
			items = append(items, *item)
		}
	}

	// finish the naming session
	mrktGroups, groups, categories := namingSession.Finish()

	return &ShopItemsResponse{items, mrktGroups, groups, categories}, nil
}

func (sic ShopItemsClient) fetchItem(
	ctx context.Context,
	asset assets.ShopAsset,
	locationId int64,
	namingSession staticdb.NamingSession,
	chnRes util.ChanSendResult[*ShopItem],
) {
	if shopPrice, err := sic.priceClient.Fetch(
		ctx,
		market.ShopParams{
			TypeId:     asset.TypeId,
			LocationId: locationId,
		},
	); err != nil { // send the error if there is one
		_ = chnRes.SendErr(err)
	} else if shopPrice == nil { // send nil if the price is nil
		_ = chnRes.SendOk(nil)
	} else { // add naming and send the item if the price is not nil
		_ = chnRes.SendOk(&ShopItem{
			Naming: namingSession.AddType(asset.TypeId),
			Asset:  asset,
			Price:  *shopPrice,
		})
	}
}
