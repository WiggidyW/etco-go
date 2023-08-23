package cancel

import (
	"context"

	"cloud.google.com/go/firestore"

	rdba "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
	rsq "github.com/WiggidyW/weve-esi/client/remotedb/rawshopqueue"
)

func CancelShopPurchase(
	rdbc *rdb.RemoteDBClient,
	ctx context.Context,
	characterId int32,
	appraisalCode string,
) error {
	fc, err := rdbc.Client(ctx)
	if err != nil {
		return err
	}
	return fc.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Remove the appraisal code from character appraisals
			if err := txDelFromCharacterShopAppraisal(
				ctx,
				tx,
				fc,
				characterId,
				// appraisalCode,
			); err != nil {
				return err
			}

			// Remove the appraisal code from shop queue
			if err := txDelFromShopQueue(
				ctx,
				tx,
				fc,
				appraisalCode,
			); err != nil {
				return err
			}

			return nil
		},
	)
}

func txDelFromShopQueue(
	ctx context.Context,
	tx *firestore.Transaction,
	fc *firestore.Client,
	appraisalCode string,
) error {
	ref := fc.Collection(rsq.COLLECTION_ID).Doc(rsq.DOC_ID)
	data := map[string]interface{}{
		rsq.FIELD_ID: firestore.ArrayRemove(appraisalCode),
	}
	return tx.Set(ref, data, firestore.MergeAll)
}

// Just sets the time of cancellation. User Appraisals includes cancelled ones.
func txDelFromCharacterShopAppraisal(
	ctx context.Context,
	tx *firestore.Transaction,
	fc *firestore.Client,
	characterId int32,
	// appraisalCode string,
) error {
	ref := rdba.CharacterRef(fc, characterId)
	data := map[string]interface{}{
		// rdba.S_CHAR_APPRAISALS:              firestore.ArrayRemove(appraisalCode),
		rdba.S_CHAR_TIME_CANCELLED_PURCHASE: firestore.ServerTimestamp,
	}
	return tx.Set(ref, data, firestore.MergeAll)
}
