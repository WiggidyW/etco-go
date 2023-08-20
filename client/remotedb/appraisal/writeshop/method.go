package writeshop

import (
	"context"

	"cloud.google.com/go/firestore"

	a "github.com/WiggidyW/weve-esi/client/appraisal"
	rdba "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
	rsq "github.com/WiggidyW/weve-esi/client/remotedb/rawshopqueue"
)

func SaveShopPurchase(
	rdbc *rdb.RemoteDBClient,
	ctx context.Context,
	appraisalCode string,
	characterId int32,
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
				characterId,
				appraisalCode,
			); err != nil {
				return err
			}

			// Append the appraisal code to shop queue
			if err := txAppendShopQueue(
				ctx,
				tx,
				fc,
				appraisalCode,
			); err != nil {
				return err
			}

			// Set the appraisal itself, with the code as the key
			if err := txSetShopAppraisal(
				ctx,
				tx,
				fc,
				characterId,
				appraisalCode,
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
		rdba.S_CHAR_APPRAISALS: firestore.ArrayUnion(appraisalCode),
	}
	return tx.Set(ref, data, firestore.MergeAll)
}

func txSetShopAppraisal(
	ctx context.Context,
	tx *firestore.Transaction,
	fc *firestore.Client,
	characterId int32,
	appraisalCode string,
	appraisal a.ShopAppraisal,
) error {
	ref := fc.Collection(rdba.SHOP_COLLECTION_ID).Doc(appraisalCode)
	data := map[string]interface{}{
		rdba.S_APPR_ITEMS:        appraisal.Items,
		rdba.S_APPR_PRICE:        appraisal.Price,
		rdba.S_APPR_TIME:         firestore.ServerTimestamp,
		rdba.S_APPR_VERSION:      appraisal.Version,
		rdba.S_APPR_LOCATION_ID:  appraisal.LocationId,
		rdba.S_APPR_CHARACTER_ID: characterId,
	}
	return tx.Set(ref, data)
}
