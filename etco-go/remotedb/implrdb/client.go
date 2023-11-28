package implrdb

import (
	"context"
)

type RemoteDBClient interface {
	ReadBuybackAppraisal(
		ctx context.Context,
		appraisalCode string,
	) (*BuybackAppraisal, error)
	ReadShopAppraisal(
		ctx context.Context,
		appraisalCode string,
	) (*ShopAppraisal, error)
	ReadUserData(
		ctx context.Context,
		characterId int32,
	) (UserData, error)
	ReadPurchaseQueue(ctx context.Context) (RawPurchaseQueue, error)
	ReadPrevContracts(ctx context.Context) (PreviousContracts, error)
	SaveBuybackAppraisal(
		ctx context.Context,
		appraisal BuybackAppraisal,
	) error
	SaveShopAppraisal(
		ctx context.Context,
		appraisal ShopAppraisal,
	) error
	DelShopPurchases(
		ctx context.Context,
		appraisalCodes ...CodeAndLocationId,
	) error
	CancelShopPurchase(
		ctx context.Context,
		characterId int32,
		appraisalCode string,
		locationId int64,
	) error
	SetPrevContracts(
		ctx context.Context,
		buybackCodes []string,
		shopCodes []string,
	) error
}
