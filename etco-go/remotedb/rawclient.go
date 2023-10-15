package remotedb

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RemoteDBClient struct {
	_client    *firestore.Client
	projectId  string
	clientOpts []option.ClientOption
	mu         *sync.Mutex
}

func NewRemoteDBClient(creds []byte, projectId string) *RemoteDBClient {
	return &RemoteDBClient{
		// _client:    nil,
		clientOpts: []option.ClientOption{
			option.WithCredentialsJSON(creds),
		},
		projectId: projectId,
		mu:        &sync.Mutex{},
	}
}

// Gets the inner client (sets it if it's nil)
func (rdbc *RemoteDBClient) innerClient() (*firestore.Client, error) {
	if rdbc._client == nil {
		// lock to prevent multiple clients from being created
		rdbc.mu.Lock()
		defer rdbc.mu.Unlock()

		// check again in case another client was created while waiting
		if rdbc._client != nil {
			return rdbc._client, nil
		}

		// create the client
		ctx := context.Background()
		var err error
		rdbc._client, err = firestore.NewClient(
			ctx,
			rdbc.projectId,
			rdbc.clientOpts...,
		)
		if err != nil {
			return nil, err
		}
	}

	return rdbc._client, nil // TODO: implement
	// panic("unimplemented")
}

func shopQueueRef(
	fc *firestore.Client,
) *firestore.DocumentRef {
	return fc.
		Collection(COLLECTION_ID_SHOP_QUEUE).
		Doc(DOCUMENT_ID_SHOP_QUEUE)
}

func userDataRef(
	characterId int32,
	fc *firestore.Client,
) *firestore.DocumentRef {
	return fc.
		Collection(COLLECTION_ID_USER_DATA).
		Doc(fmt.Sprintf("%d", characterId))
}

func buybackAppraisalRef(
	appraisalCode string,
	fc *firestore.Client,
) *firestore.DocumentRef {
	return fc.
		Collection(COLLECTION_ID_BUYBACK_APPRAISALS).
		Doc(appraisalCode)
}

func shopAppraisalRef(
	appraisalCode string,
	fc *firestore.Client,
) *firestore.DocumentRef {
	return fc.
		Collection(COLLECTION_ID_SHOP_APPRAISALS).
		Doc(appraisalCode)
}

func txSetShopQueue(
	fc *firestore.Client,
	tx *firestore.Transaction,
	data map[string]interface{},
	opts ...firestore.SetOption,
) error {
	return tx.Set(shopQueueRef(fc), data, opts...)
}

func txSetUserData(
	characterId int32,
	fc *firestore.Client,
	tx *firestore.Transaction,
	data map[string]interface{},
	opts ...firestore.SetOption,
) error {
	return tx.Set(userDataRef(characterId, fc), data, opts...)
}

func txSetBuybackAppraisal(
	appraisalCode string,
	fc *firestore.Client,
	tx *firestore.Transaction,
	data map[string]interface{},
	opts ...firestore.SetOption,
) error {
	return tx.Set(buybackAppraisalRef(appraisalCode, fc), data, opts...)
}

func txSetShopAppraisal(
	appraisalCode string,
	fc *firestore.Client,
	tx *firestore.Transaction,
	data map[string]interface{},
	opts ...firestore.SetOption,
) error {
	return tx.Set(shopAppraisalRef(appraisalCode, fc), data, opts...)
}

func txDataRemoveManyFromShopQueue(
	remove ...string,
) (map[string]interface{}, firestore.SetOption) {
	removeAsAny := make([]any, len(remove))
	for i, v := range remove {
		removeAsAny[i] = v
	}

	return map[string]interface{}{
		FIELD_SHOP_QUEUE_SHOP_QUEUE: firestore.ArrayRemove(
			removeAsAny...,
		),
	}, firestore.MergeAll
}

func txDataRemoveOneFromShopQueue(
	remove string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_SHOP_QUEUE_SHOP_QUEUE: firestore.ArrayRemove(remove),
	}, firestore.MergeAll
}

func txDataAppendToShopQueue(
	append string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_SHOP_QUEUE_SHOP_QUEUE: firestore.ArrayUnion(append),
	}, firestore.MergeAll
}

func txDataRemovePurchaseFromUserData() (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_USER_DATA_TIME_CANCELLED_PURCHASE: firestore.
			ServerTimestamp,
	}, firestore.MergeAll
}

func txDataAppendPurchaseToUserData(
	appraisalCode string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_USER_DATA_SHOP_APPRAISALS: firestore.ArrayUnion(
			appraisalCode,
		),
		FIELD_USER_DATA_TIME_MADE_PURCHASE: firestore.ServerTimestamp,
	}, firestore.MergeAll
}

func txDataAppendBuybackAppraisalToUserData(
	appraisalCode string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_USER_DATA_BUYBACK_APPRAISALS: firestore.ArrayUnion(
			appraisalCode,
		),
	}, firestore.MergeAll
}

func txDataSetShopAppraisal(
	appraisal ShopAppraisal,
) map[string]interface{} {
	return map[string]interface{}{
		FIELD_SHOP_APPRAISAL_TIME:         firestore.ServerTimestamp,
		FIELD_SHOP_APPRAISAL_ITEMS:        appraisal.Items,
		FIELD_SHOP_APPRAISAL_PRICE:        appraisal.Price,
		FIELD_SHOP_APPRAISAL_TAX_RATE:     appraisal.TaxRate,
		FIELD_SHOP_APPRAISAL_TAX:          appraisal.Tax,
		FIELD_SHOP_APPRAISAL_VERSION:      appraisal.Version,
		FIELD_SHOP_APPRAISAL_CHARACTER_ID: appraisal.CharacterId,
		FIELD_SHOP_APPRAISAL_LOCATION_ID:  appraisal.LocationId,
	}
}

func txDataSetBuybackAppraisal(
	appraisal BuybackAppraisal,
) map[string]interface{} {
	return map[string]interface{}{
		FIELD_BUYBACK_APPRAISAL_TIME:         firestore.ServerTimestamp,
		FIELD_BUYBACK_APPRAISAL_ITEMS:        appraisal.Items,
		FIELD_BUYBACK_APPRAISAL_PRICE:        appraisal.Price,
		FIELD_BUYBACK_APPRAISAL_TAX_RATE:     appraisal.TaxRate,
		FIELD_BUYBACK_APPRAISAL_TAX:          appraisal.Tax,
		FIELD_BUYBACK_APPRAISAL_VERSION:      appraisal.Version,
		FIELD_BUYBACK_APPRAISAL_CHARACTER_ID: appraisal.CharacterId,
		FIELD_BUYBACK_APPRAISAL_SYSTEM_ID:    appraisal.SystemId,
		FIELD_BUYBACK_APPRAISAL_FEE:          appraisal.Fee,
		FIELD_BUYBACK_APPRAISAL_FEE_PER_M3:   appraisal.FeePerM3,
	}
}

func Read[V any](
	ctx context.Context,
	fc *firestore.Client,
	ref *firestore.DocumentRef,
	val *V,
) (exists bool, err error) {
	doc, err := ref.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	if err := doc.DataTo(val); err != nil {
		return true, err
	} else {
		return true, nil
	}
}

func (rdbc *RemoteDBClient) ReadShopQueue(
	ctx context.Context,
) (val ShopQueue, err error) {
	if fc, err := rdbc.innerClient(); err != nil {
		return val, err
	} else {
		ref := shopQueueRef(fc)
		_, err = Read(ctx, fc, ref, &val) // ignore exists false and return empty
		return val, err
	}
}

func (rdbc *RemoteDBClient) ReadUserData(
	ctx context.Context,
	characterId int32,
) (val UserData, err error) {
	if fc, err := rdbc.innerClient(); err != nil {
		return val, err
	} else {
		ref := userDataRef(characterId, fc)
		_, err = Read(ctx, fc, ref, &val) // ignore exists false and return empty
		return val, err
	}
}

func (rdbc *RemoteDBClient) ReadShopAppraisal(
	ctx context.Context,
	appraisalCode string,
) (exists bool, val ShopAppraisal, err error) {
	if fc, err := rdbc.innerClient(); err != nil {
		return false, val, err
	} else {
		ref := shopAppraisalRef(appraisalCode, fc)
		exists, err = Read(ctx, fc, ref, &val)
		if exists && err == nil {
			val.Code = appraisalCode
		}
		return exists, val, err
	}
}

func (rdbc *RemoteDBClient) ReadBuybackAppraisal(
	ctx context.Context,
	appraisalCode string,
) (exists bool, val BuybackAppraisal, err error) {
	if fc, err := rdbc.innerClient(); err != nil {
		return false, val, err
	} else {
		ref := buybackAppraisalRef(appraisalCode, fc)
		exists, err = Read(ctx, fc, ref, &val)
		if exists && err == nil {
			val.Code = appraisalCode
		}
		return exists, val, err
	}
}

func (rdbc *RemoteDBClient) SaveShopPurchase(
	ctx context.Context,
	appraisal ShopAppraisal,
) error {
	fc, err := rdbc.innerClient()
	if err != nil {
		return err
	}
	return fc.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Append the purchase to user data
			txudData, txudOpts := txDataAppendPurchaseToUserData(
				appraisal.Code,
			)
			if err := txSetUserData(
				appraisal.CharacterId,
				fc,
				tx,
				txudData,
				txudOpts,
			); err != nil {
				return err
			}

			// Append the appraisal code to shop queue
			txsqData, txsqOpts := txDataAppendToShopQueue(
				appraisal.Code,
			)
			if err := txSetShopQueue(
				fc,
				tx,
				txsqData,
				txsqOpts,
			); err != nil {
				return err
			}

			// Set the appraisal itself, with the code as the key
			txsaData := txDataSetShopAppraisal(appraisal)
			if err := txSetShopAppraisal(
				appraisal.Code,
				fc,
				tx,
				txsaData,
			); err != nil {
				return err
			}

			return nil
		},
	)
}

func (rdbc *RemoteDBClient) DelShopPurchases(
	ctx context.Context,
	appraisalCodes ...string,
) error {
	fc, err := rdbc.innerClient()
	if err != nil {
		return err
	}
	return fc.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Remove the appraisal codes from shop queue
			txsqData, txsqOpts := txDataRemoveManyFromShopQueue(
				appraisalCodes...,
			)
			if err := txSetShopQueue(
				fc,
				tx,
				txsqData,
				txsqOpts,
			); err != nil {
				return err
			}

			return nil
		},
	)
}

func (rdbc *RemoteDBClient) CancelShopPurchase(
	ctx context.Context,
	characterId int32,
	appraisalCode string,
) error {
	fc, err := rdbc.innerClient()
	if err != nil {
		return err
	}
	return fc.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Remove the appraisal code from user data
			txudData, txudOpts := txDataRemovePurchaseFromUserData()
			if err := txSetUserData(
				characterId,
				fc,
				tx,
				txudData,
				txudOpts,
			); err != nil {
				return err
			}

			// Remove the appraisal code from shop queue
			txsqData, txsqOpts := txDataRemoveOneFromShopQueue(
				appraisalCode,
			)
			if err := txSetShopQueue(
				fc,
				tx,
				txsqData,
				txsqOpts,
			); err != nil {
				return err
			}

			return nil
		},
	)
}

func (rdbc *RemoteDBClient) SaveBuybackAppraisal(
	ctx context.Context,
	appraisal BuybackAppraisal,
) error {
	fc, err := rdbc.innerClient()
	if err != nil {
		return err
	}
	return fc.RunTransaction(
		ctx,
		func(ctx context.Context, tx *firestore.Transaction) error {
			// Append the appraisal code to character appraisals
			if appraisal.CharacterId != nil {
				txudData, txudOpts := txDataAppendBuybackAppraisalToUserData(
					appraisal.Code,
				)
				if err := txSetUserData(
					*appraisal.CharacterId,
					fc,
					tx,
					txudData,
					txudOpts,
				); err != nil {
					return err
				}
			}

			// Set the appraisal itself, with the code as the key
			txbaData := txDataSetBuybackAppraisal(appraisal)
			if err := txSetBuybackAppraisal(
				appraisal.Code,
				fc,
				tx,
				txbaData,
			); err != nil {
				return err
			}

			return nil
		},
	)
}
