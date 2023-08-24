package writebuyback

import (
	"context"

	"cloud.google.com/go/firestore"

	a "github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	rdba "github.com/WiggidyW/eve-trading-co-go/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/eve-trading-co-go/client/remotedb/internal"
)

func SaveBuybackAppraisal(
	rdbc *rdb.RemoteDBClient,
	ctx context.Context,
	appraisal a.BuybackAppraisal,
) error {
	fc, err := rdbc.Client(ctx)
	if err != nil {
		return err
	}
	return fc.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Append the appraisal code to character appraisals
			if appraisal.CharacterId != nil {
				if err := txAppendCharacterBuybackAppraisal(
					ctx,
					tx,
					fc,
					*appraisal.CharacterId,
					appraisal.Code,
				); err != nil {
					return err
				}
			}

			// Set the appraisal itself, with the code as the key
			if err := txSetBuybackAppraisal(
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

func txAppendCharacterBuybackAppraisal(
	ctx context.Context,
	tx *firestore.Transaction,
	fc *firestore.Client,
	characterId int32,
	appraisalCode string,
) error {
	ref := rdba.CharacterRef(fc, characterId)
	data := map[string]interface{}{
		rdba.B_CHAR_APPRAISALS: firestore.ArrayUnion(appraisalCode),
	}
	return tx.Set(ref, data, firestore.MergeAll)
}

func txSetBuybackAppraisal(
	ctx context.Context,
	tx *firestore.Transaction,
	fc *firestore.Client,
	appraisal a.BuybackAppraisal,
) error {
	ref := fc.Collection(rdba.BUYBACK_COLLECTION_ID).Doc(appraisal.Code)
	data := map[string]interface{}{
		rdba.B_APPR_ITEMS:        appraisal.Items,
		rdba.B_APPR_PRICE:        appraisal.Price,
		rdba.B_APPR_TIME:         firestore.ServerTimestamp,
		rdba.B_APPR_VERSION:      appraisal.Version,
		rdba.B_APPR_SYSTEM_ID:    appraisal.SystemId,
		rdba.B_APPR_CHARACTER_ID: appraisal.CharacterId,
	}
	return tx.Set(ref, data)
}
