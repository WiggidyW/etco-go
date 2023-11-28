package remotedb

import (
	"context"
	"fmt"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/remotedb/firestorerdb"
	"github.com/WiggidyW/etco-go/remotedb/implrdb"
	"github.com/WiggidyW/etco-go/remotedb/mysqlrdb"
)

var client implrdb.RemoteDBClient

func init() {
	switch build.REMOTEDB {
	case build.RDBFirestore:
		client = firestorerdb.NewFSClient()
	case build.RDBMySQL:
		client = mysqlrdb.NewMySQLClient()
	default:
		panic(fmt.Sprintf("invalid REMOTEDB: %d", build.REMOTEDB))
	}
}

type CodeAndLocationId = implrdb.CodeAndLocationId
type ICodeAndLocationId interface {
	GetCode() string
	GetLocationId() int64
}

func saveBuybackAppraisal(ctx context.Context, appraisal BuybackAppraisal) error {
	return client.SaveBuybackAppraisal(ctx, appraisal)
}

func saveShopAppraisal(ctx context.Context, appraisal ShopAppraisal) error {
	return client.SaveShopAppraisal(ctx, appraisal)
}

func cancelShopPurchase(
	ctx context.Context,
	characterId int32,
	appraisalCode string,
	locationId int64,
) error {
	return client.CancelShopPurchase(ctx, characterId, appraisalCode, locationId)
}

func setPrevContracts(
	ctx context.Context,
	buybackCodes []string,
	shopCodes []string,
) error {
	return client.SetPrevContracts(ctx, buybackCodes, shopCodes)
}

func readBuybackAppraisal(
	ctx context.Context,
	appraisalCode string,
) (*BuybackAppraisal, error) {
	return client.ReadBuybackAppraisal(ctx, appraisalCode)
}

func readShopAppraisal(
	ctx context.Context,
	appraisalCode string,
) (*ShopAppraisal, error) {
	return client.ReadShopAppraisal(ctx, appraisalCode)
}

func readUserData(
	ctx context.Context,
	characterId int32,
) (UserData, error) {
	return client.ReadUserData(ctx, characterId)
}

func readPurchaseQueue(ctx context.Context) (RawPurchaseQueue, error) {
	return client.ReadPurchaseQueue(ctx)
}

func readPrevContracts(ctx context.Context) (PreviousContracts, error) {
	return client.ReadPrevContracts(ctx)
}

func delShopPurchases[C ICodeAndLocationId](
	ctx context.Context,
	icodes ...C,
) error {
	codes := make([]CodeAndLocationId, len(icodes))
	for i, icode := range icodes {
		codes[i] = CodeAndLocationId{
			Code:       icode.GetCode(),
			LocationId: icode.GetLocationId(),
		}
	}
	return client.DelShopPurchases(ctx, codes...)
}
