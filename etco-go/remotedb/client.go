package remotedb

import (
	"context"
	"fmt"
	"sync"

	build "github.com/WiggidyW/etco-go/buildconstants"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	client *fsClient
)

func init() {
	client = newFSClient(
		[]byte(build.REMOTEDB_CREDS_JSON),
		build.REMOTEDB_PROJECT_ID,
	)
}

type fsClient struct {
	_client    *firestore.Client
	projectId  string
	clientOpts []option.ClientOption
	mu         *sync.Mutex
}

func newFSClient(creds []byte, projectId string) *fsClient {
	return &fsClient{
		// _client:    nil,
		clientOpts: []option.ClientOption{
			option.WithCredentialsJSON(creds),
		},
		projectId: projectId,
		mu:        &sync.Mutex{},
	}
}

// Gets the inner client (sets it if it's nil)
func (c *fsClient) innerClient() (*firestore.Client, error) {
	if c._client == nil {
		// lock to prevent multiple clients from being created
		c.mu.Lock()
		defer c.mu.Unlock()

		// check again in case another client was created while waiting
		if c._client != nil {
			return c._client, nil
		}

		// create the client
		ctx := context.Background()
		var err error
		c._client, err = firestore.NewClient(
			ctx,
			c.projectId,
			c.clientOpts...,
		)
		if err != nil {
			return nil, err
		}
	}

	return c._client, nil // TODO: implement
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

func txSetPurchaseQueue(
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

func txDataRemoveManyFromPurchaseQueue(
	remove ...string,
) (map[string]interface{}, firestore.SetOption) {
	removeAsAny := make([]any, len(remove))
	for i, v := range remove {
		removeAsAny[i] = v
	}

	return map[string]interface{}{
		FIELD_SHOP_QUEUE_SHOP_QUEUE: firestore.ArrayRemove(removeAsAny...),
	}, firestore.MergeAll
}

func txDataRemoveOneFromPurchaseQueue(
	remove string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_SHOP_QUEUE_SHOP_QUEUE: firestore.ArrayRemove(remove),
	}, firestore.MergeAll
}

func txDataAppendToPurchaseQueue(
	append string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_SHOP_QUEUE_SHOP_QUEUE: firestore.ArrayUnion(append),
	}, firestore.MergeAll
}

func txDataCancelPurchaseUserData() (
	map[string]interface{},
	firestore.SetOption,
) {
	return map[string]interface{}{
		FIELD_USER_DATA_TIME_CANCELLED_PURCHASE: firestore.ServerTimestamp,
	}, firestore.MergeAll
}

func txDataAppendShopAppraisalToUserData(
	appraisalCode string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_USER_DATA_SHOP_APPRAISALS:    firestore.ArrayUnion(appraisalCode),
		FIELD_USER_DATA_TIME_MADE_PURCHASE: firestore.ServerTimestamp,
	}, firestore.MergeAll
}

func txDataAppendBuybackAppraisalToUserData(
	appraisalCode string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_USER_DATA_BUYBACK_APPRAISALS: firestore.ArrayUnion(appraisalCode),
	}, firestore.MergeAll
}

func txDataSetShopAppraisal(
	appraisal ShopAppraisal,
) map[string]interface{} {
	return map[string]interface{}{
		FIELD_SHOP_APPRAISAL_TIME:         appraisal.Time,
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
		FIELD_BUYBACK_APPRAISAL_TIME:         appraisal.Time,
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

func read[V any](
	ctx context.Context,
	fc *firestore.Client,
	ref *firestore.DocumentRef,
) (val *V, err error) {
	doc, err := ref.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	val = new(V)
	err = doc.DataTo(val)
	if err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

func (c *fsClient) readPurchaseQueue(
	ctx context.Context,
) (*PurchaseQueue, error) {
	fc, err := c.innerClient()
	if err != nil {
		return nil, err
	}

	ref := shopQueueRef(fc)
	return read[PurchaseQueue](ctx, fc, ref)
}

func (c *fsClient) readUserData(
	ctx context.Context,
	characterId int32,
) (*UserData, error) {
	fc, err := c.innerClient()
	if err != nil {
		return nil, err
	}

	ref := userDataRef(characterId, fc)
	return read[UserData](ctx, fc, ref)
}

func (c *fsClient) readShopAppraisal(
	ctx context.Context,
	appraisalCode string,
) (*ShopAppraisal, error) {
	fc, err := c.innerClient()
	if err != nil {
		return nil, err
	}

	ref := shopAppraisalRef(appraisalCode, fc)
	rep, err := read[ShopAppraisal](ctx, fc, ref)
	if rep != nil {
		rep.Code = appraisalCode
	}
	return rep, err
}

func (c *fsClient) readBuybackAppraisal(
	ctx context.Context,
	appraisalCode string,
) (*BuybackAppraisal, error) {
	fc, err := c.innerClient()
	if err != nil {
		return nil, err
	}

	ref := buybackAppraisalRef(appraisalCode, fc)
	rep, err := read[BuybackAppraisal](ctx, fc, ref)
	if rep != nil {
		rep.Code = appraisalCode
	}
	return rep, err
}

func (c *fsClient) saveShopAppraisal(
	ctx context.Context,
	appraisal ShopAppraisal,
) error {
	fc, err := c.innerClient()
	if err != nil {
		return err
	}

	txFunc := func(ctx context.Context, tx *firestore.Transaction) error {
		// Append the appraisal to user data
		txudData, txudOpts := txDataAppendShopAppraisalToUserData(appraisal.Code)
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
		txsqData, txsqOpts := txDataAppendToPurchaseQueue(appraisal.Code)
		if err := txSetPurchaseQueue(
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
	}
	return fc.RunTransaction(ctx, txFunc)
}

func (c *fsClient) delShopPurchases(
	ctx context.Context,
	appraisalCodes ...string,
) error {
	fc, err := c.innerClient()
	if err != nil {
		return err
	}

	txFunc := func(ctx context.Context, tx *firestore.Transaction) error {
		// Remove the appraisal codes from shop queue
		txsqData, txsqOpts := txDataRemoveManyFromPurchaseQueue(appraisalCodes...)
		if err := txSetPurchaseQueue(
			fc,
			tx,
			txsqData,
			txsqOpts,
		); err != nil {
			return err
		}

		return nil
	}
	return fc.RunTransaction(ctx, txFunc)
}

func (c *fsClient) cancelShopPurchase(
	ctx context.Context,
	characterId int32,
	appraisalCode string,
) error {
	fc, err := c.innerClient()
	if err != nil {
		return err
	}

	txFunc := func(ctx context.Context, tx *firestore.Transaction) error {
		// Set cancellation time in user data
		txudData, txudOpts := txDataCancelPurchaseUserData()
		if err := txSetUserData(
			characterId,
			fc,
			tx,
			txudData,
			txudOpts,
		); err != nil {
			return err
		}

		// Remove the appraisal code from purchase queue
		txsqData, txsqOpts := txDataRemoveOneFromPurchaseQueue(appraisalCode)
		if err := txSetPurchaseQueue(
			fc,
			tx,
			txsqData,
			txsqOpts,
		); err != nil {
			return err
		}

		return nil
	}
	return fc.RunTransaction(ctx, txFunc)
}

func (c *fsClient) saveBuybackAppraisal(
	ctx context.Context,
	appraisal BuybackAppraisal,
) error {
	fc, err := c.innerClient()
	if err != nil {
		return err
	}

	txFunc := func(ctx context.Context, tx *firestore.Transaction) error {
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
	}
	return fc.RunTransaction(ctx, txFunc)
}
