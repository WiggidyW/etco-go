package write

import (
	"context"

	"cloud.google.com/go/firestore"

	a "github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	rdba "github.com/WiggidyW/eve-trading-co-go/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/eve-trading-co-go/client/remotedb/internal"
	rsq "github.com/WiggidyW/eve-trading-co-go/client/remotedb/rawshopqueue"
)

func SaveShopPurchase(
	rdbc *rdb.RemoteDBClient,
	ctx context.Context,
	appraisal a.ShopAppraisal,
) error {
	fc, err := rdbc.Client(ctx)
	if err != nil {
		return err
	}
	return fc.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Append the appraisal code to character appraisals
			if err := txAppendCharacterShopAppraisal(
				ctx,
				tx,
				fc,
				appraisal.CharacterId,
				appraisal.Code,
			); err != nil {
				return err
			}

			// Append the appraisal code to shop queue
			if err := txAppendShopQueue(
				ctx,
				tx,
				fc,
				appraisal.Code,
			); err != nil {
				return err
			}

			// Set the appraisal itself, with the code as the key
			if err := txSetShopAppraisal(
				ctx,
				tx,
				fc,
				appraisal,
			); err != nil {
				return err
			}

			return nil
		},
	)
}

func txAppendShopQueue(
	ctx context.Context,
	tx *firestore.Transaction,
	fc *firestore.Client,
	appraisalCode string,
) error {
	ref := fc.Collection(rsq.COLLECTION_ID).Doc(rsq.DOC_ID)
	data := map[string]interface{}{
		rsq.FIELD_ID: firestore.ArrayUnion(appraisalCode),
	}
	return tx.Set(ref, data, firestore.MergeAll)
}

func txAppendCharacterShopAppraisal(
	ctx context.Context,
	tx *firestore.Transaction,
	fc *firestore.Client,
	characterId int32,
	appraisalCode string,
) error {
	ref := rdba.CharacterRef(fc, characterId)
	data := map[string]interface{}{
		rdba.S_CHAR_APPRAISALS:         firestore.ArrayUnion(appraisalCode),
		rdba.S_CHAR_TIME_MADE_PURCHASE: firestore.ServerTimestamp,
	}
	return tx.Set(ref, data, firestore.MergeAll)
}

func txSetShopAppraisal(
	ctx context.Context,
	tx *firestore.Transaction,
	fc *firestore.Client,
	appraisal a.ShopAppraisal,
) error {
	ref := fc.Collection(rdba.SHOP_COLLECTION_ID).Doc(appraisal.Code)
	data := map[string]interface{}{
		rdba.S_APPR_ITEMS:        appraisal.Items,
		rdba.S_APPR_PRICE:        appraisal.Price,
		rdba.S_APPR_TIME:         firestore.ServerTimestamp,
		rdba.S_APPR_VERSION:      appraisal.Version,
		rdba.S_APPR_LOCATION_ID:  appraisal.LocationId,
		rdba.S_APPR_CHARACTER_ID: appraisal.CharacterId,
	}
	return tx.Set(ref, data)
}
