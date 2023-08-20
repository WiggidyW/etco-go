package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CHARACTERS_COLLECTION_ID string = "users"
	// CHARACTER_SHOP_CODES_FIELD      string = "ShopCodes"
	// CHARACTER_BUYBACK_CODES_FIELD   string = "BuybackCodes"
	SHOP_QUEUE_COLLECTION_ID string = "shop_queue"
	SHOP_QUEUE_DOCUMENT_ID   string = "shop_queue"
	// SHOP_QUEUE_CODES_FIELD          string = "ShopCodes"
	SHOP_APPRAISAL_COLLECTION_ID    string = "shop_appraisals"
	BUYBACK_APPRAISAL_COLLECTION_ID string = "buyback_appraisals"
)

type DatabaseClient struct {
	firestoreClient *firestore.Client
}

func (dbc *DatabaseClient) GetCharacterCodes(
	ctx context.Context,
	characterID int32,
) (*CharacterCodes, error) {
	// character codes get parameters
	characterCodesRef := dbc.firestoreClient.
		Collection(CHARACTERS_COLLECTION_ID).
		Doc(fmt.Sprintf("%d", characterID))

	// get the character codes document
	characterCodesDoc, err := characterCodesRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// return an empty character codes
			return &CharacterCodes{}, nil
		} else {
			return nil, err
		}
	}

	// parse the character codes document
	cc := new(CharacterCodes)
	if err := characterCodesDoc.DataTo(cc); err != nil {
		return nil, err
	} else {
		// return the character codes
		return cc, nil
	}
}

func (dbc *DatabaseClient) GetShopQueue(
	ctx context.Context,
) ([]string, error) {
	// shop queue get parameters
	shopQueueRef := dbc.firestoreClient.
		Collection(SHOP_QUEUE_COLLECTION_ID).
		Doc(SHOP_QUEUE_DOCUMENT_ID)

	// get the shop queue document
	shopQueueDoc, err := shopQueueRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// return an empty shop queue
			return []string{}, nil
		} else {
			return nil, err
		}
	}

	// parse the shop queue document
	sq := make([]string, 0)
	if err := shopQueueDoc.DataTo(&sq); err != nil {
		return nil, err
	} else {
		// return the shop queue
		return sq, nil
	}
}

func (dbc *DatabaseClient) GetBuybackAppraisal(
	ctx context.Context,
	appraisalKey string,
) (ba BuybackAppraisal, err error) {
	err = dbc.getAppraisal(
		ctx,
		appraisalKey,
		BUYBACK_APPRAISAL_COLLECTION_ID,
		&ba,
	)
	return ba, err
}

func (dbc *DatabaseClient) GetShopAppraisal(
	ctx context.Context,
	appraisalKey string,
) (sa ShopAppraisal, err error) {
	err = dbc.getAppraisal(
		ctx,
		appraisalKey,
		SHOP_APPRAISAL_COLLECTION_ID,
		&sa,
	)
	return sa, err
}

func (dbc *DatabaseClient) getAppraisal(
	ctx context.Context,
	appraisalKey string,
	collectionId string,
	dataTo interface{}, // must be a pointer
) error {
	// appraisal get parameters
	appraisalRef := dbc.firestoreClient.
		Collection(collectionId).
		Doc(appraisalKey)

	// get the appraisal document
	appraisalDoc, err := appraisalRef.Get(ctx)
	if err != nil {
		return err
	}

	// parse the appraisal document
	if err := appraisalDoc.DataTo(dataTo); err != nil {
		return err
	}

	return nil
}

func (dbc *DatabaseClient) CancelPurchase(
	ctx context.Context,
	characterId int32,
	appraisalKey string,
) error {
	// shop queue del parameters
	shopQueueRef := dbc.firestoreClient.
		Collection(SHOP_QUEUE_COLLECTION_ID).
		Doc(SHOP_QUEUE_DOCUMENT_ID)
	shopQueueData := map[string]interface{}{
		"ShopCodes": firestore.ArrayRemove(appraisalKey),
	}

	// run the transaction
	return dbc.firestoreClient.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			return tx.Set(
				shopQueueRef,
				shopQueueData,
				firestore.MergeAll,
			)
		},
	)
}

func SaveBuybackAppraisal[
	B IBuybackAppraisal[I],
	I IBuybackAppraisalItem[CI],
	CI IBuybackAppraisalChildItem,
](
	dbc *DatabaseClient,
	ctx context.Context,
	characterId int32,
	appraisalKey string,
	appraisal B,
) error {
	// // helper types for firestore
	type saveBuybackChildItem struct {
		TypeId       int32
		Quantity     float64
		PricePerUnit float64
		Description  string
	}

	type saveBuybackItem struct {
		TypeId       int32
		Quantity     int64
		PricePerUnit float64
		Description  string
		Children     []saveBuybackChildItem
	}
	//

	// character appraisals append parameters
	characterRef := dbc.firestoreClient.
		Collection(CHARACTERS_COLLECTION_ID).
		Doc(fmt.Sprintf("%d", characterId))
	characterData := map[string]interface{}{
		"BuybackCodes": firestore.ArrayUnion(
			appraisalKey,
		),
	}

	// // appraisal set parameters
	appraisalRef := dbc.firestoreClient.
		Collection(BUYBACK_APPRAISAL_COLLECTION_ID).
		Doc(appraisalKey)

	// create the items field
	aItems := appraisal.GetItems()
	aDataItems := make([]saveBuybackItem, 0, len(aItems))
	for _, appraisalItem := range aItems {

		// create the children field
		acItems := appraisalItem.GetChildren()
		acDataItems := make([]saveBuybackChildItem, 0, len(acItems))
		for _, acItem := range acItems {
			acDataItems = append(acDataItems, saveBuybackChildItem{
				TypeId:       acItem.GetTypeId(),
				Quantity:     acItem.GetQuantity(),
				PricePerUnit: acItem.GetPricePerUnit(),
				Description:  acItem.GetDescription(),
			})
		}

		// append the item
		aDataItems = append(aDataItems, saveBuybackItem{
			TypeId:       appraisalItem.GetTypeId(),
			Quantity:     appraisalItem.GetQuantity(),
			PricePerUnit: appraisalItem.GetPricePerUnit(),
			Description:  appraisalItem.GetDescription(),
			Children:     acDataItems,
		})
	}

	// create the appraisal data
	appraisalData := newAppraisalData[B, I, saveBuybackItem, int32](
		characterId,
		appraisal,
		aDataItems,
		"SystemId",
		appraisal.GetSystemId(),
	)
	// //

	// run the transaction
	return dbc.firestoreClient.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Append the appraisal key to character appraisals
			if err := tx.Set(
				characterRef,
				characterData,
			); err != nil {
				return err
			}

			// Set the appraisal
			return tx.Set(
				appraisalRef,
				appraisalData,
			)
		},
	)
}

func MakePurchase[S IShopAppraisal[I], I IShopAppraisalItem](
	dbc *DatabaseClient,
	ctx context.Context,
	characterId int32,
	appraisalKey string,
	appraisal S,
) error {
	// helper type for firestore
	type makePurchaseItem struct {
		TypeId       int32
		Quantity     int64
		PricePerUnit float64
		Description  string
	}

	// character appraisals append parameters
	characterRef := dbc.firestoreClient.
		Collection(CHARACTERS_COLLECTION_ID).
		Doc(fmt.Sprintf("%d", characterId))
	characterData := map[string]interface{}{
		"ShopCodes": firestore.ArrayUnion(appraisalKey),
	}

	// shop queue append parameters
	shopQueueRef := dbc.firestoreClient.
		Collection(SHOP_QUEUE_COLLECTION_ID).
		Doc(SHOP_QUEUE_DOCUMENT_ID)
	shopQueueData := map[string]interface{}{
		"ShopCodes": firestore.ArrayUnion(appraisalKey),
	}

	// // appraisal set parameters
	appraisalRef := dbc.firestoreClient.
		Collection(SHOP_APPRAISAL_COLLECTION_ID).
		Doc(appraisalKey)

	// create the items field
	aItems := appraisal.GetItems()
	aDataItems := make([]makePurchaseItem, 0, len(aItems))
	for _, item := range aItems {
		// append the item
		aDataItems = append(aDataItems, makePurchaseItem{
			TypeId:       item.GetTypeId(),
			Quantity:     item.GetQuantity(),
			PricePerUnit: item.GetPricePerUnit(),
			Description:  item.GetDescription(),
		})
	}

	// create the appraisal data
	appraisalData := newAppraisalData[S, I, makePurchaseItem, int64](
		characterId,
		appraisal,
		aDataItems,
		"LocationId",
		appraisal.GetLocationId(),
	)
	// //

	// run the transaction
	return dbc.firestoreClient.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Append the appraisal key to character appraisals
			if err := tx.Set(
				characterRef,
				characterData,
				firestore.MergeAll,
			); err != nil {
				return err
			}

			// Append the appraisal key to shop queue
			if err := tx.Set(
				shopQueueRef,
				shopQueueData,
				firestore.MergeAll,
			); err != nil {
				return err
			}

			// Set the appraisal
			return tx.Set(
				appraisalRef,
				appraisalData,
			)
		},
	)
}

func newAppraisalData[A IAppraisal[AI], AI any, I any, L any](
	characterId int32,
	appraisal A,
	items []I,
	locationKey string,
	locationVal L,
) map[string]interface{} {
	return map[string]interface{}{
		"Items":       items,
		"Price":       appraisal.GetPrice(),
		"Time":        firestore.ServerTimestamp, // firestore will set it for us
		"Version":     appraisal.GetVersion(),
		locationKey:   locationVal,
		"CharacterId": characterId,
	}
}
