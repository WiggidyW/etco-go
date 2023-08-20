package writebuyback

import (
	"context"

	"cloud.google.com/go/firestore"

	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

func SaveBuybackAppraisal[
	B a.IBuybackAppraisal[I],
	I a.IBuybackParentItem[CI],
	CI a.IBuybackChildItem,
](
	rdbc *rdb.RemoteDBClient,
	ctx context.Context,
	appraisalCode string,
	characterId *int32,
	iAppraisal B,
) error {
	fc, err := rdbc.Client(ctx)
	if err != nil {
		return err
	}
	return fc.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Append the appraisal code to character appraisals
			if characterId != nil {
				if err := txAppendCharacterBuybackAppraisal(
					ctx,
					tx,
					fc,
					*characterId,
					appraisalCode,
				); err != nil {
					return err
				}
			}

			// Set the appraisal itself, with the code as the key
			if err := txSetBuybackAppraisal[B, I, CI](
				ctx,
				tx,
				fc,
				characterId,
				appraisalCode,
				iAppraisal,
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
	ref := a.CharacterRef(fc, characterId)
	data := map[string]interface{}{
		a.B_CHAR_APPRAISALS: firestore.ArrayUnion(appraisalCode),
	}
	return tx.Set(ref, data, firestore.MergeAll)
}

func txSetBuybackAppraisal[
	B a.IBuybackAppraisal[I],
	I a.IBuybackParentItem[CI],
	CI a.IBuybackChildItem,
](
	ctx context.Context,
	tx *firestore.Transaction,
	fc *firestore.Client,
	characterId *int32,
	appraisalCode string,
	iAppraisal B,
) error {
	ref := fc.Collection(a.BUYBACK_COLLECTION_ID).Doc(appraisalCode)
	data := map[string]interface{}{
		a.B_APPR_ITEMS: a.NewBuybackParentItems[I, CI](
			iAppraisal.GetItems(),
		),
		a.B_APPR_PRICE:        iAppraisal.GetPrice(),
		a.B_APPR_TIME:         firestore.ServerTimestamp,
		a.B_APPR_VERSION:      iAppraisal.GetVersion(),
		a.B_APPR_LOCATION_ID:  iAppraisal.GetLocationId(),
		a.B_APPR_CHARACTER_ID: characterId,
	}
	return tx.Set(ref, data)
}
