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

type ICodeAndLocationId interface {
	GetCode() string
	GetLocationId() int64
}

type CodeAndLocationId struct {
	Code       string
	LocationId int64
}

func NewCodeAndLocationId(code string, locationId int64) CodeAndLocationId {
	return CodeAndLocationId{Code: code, LocationId: locationId}
}

func (c CodeAndLocationId) GetCode() string      { return c.Code }
func (c CodeAndLocationId) GetLocationId() int64 { return c.LocationId }

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

func prevContractsRef(
	fc *firestore.Client,
) *firestore.DocumentRef {
	return fc.
		Collection(COLLECTION_ID_CONTRACTS).
		Doc(DOCUMENT_ID_PREV_CONTRACTS)
}

func purchaseQueueRef(
	fc *firestore.Client,
) *firestore.DocumentRef {
	return fc.
		Collection(COLLECTION_ID_PURCHASE_QUEUE).
		Doc(DOCUMENT_ID_PURCHASE_QUEUE)
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

func txSetPrevContracts(
	fc *firestore.Client,
	tx *firestore.Transaction,
	data map[string]interface{},
	opts ...firestore.SetOption,
) error {
	return tx.Set(prevContractsRef(fc), data, opts...)
}

func txSetPurchaseQueue(
	fc *firestore.Client,
	tx *firestore.Transaction,
	data map[string]interface{},
	opts ...firestore.SetOption,
) error {
	return tx.Set(purchaseQueueRef(fc), data, opts...)
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

func txDataSetPrevContracts(
	buybackCodes []string,
	shopCodes []string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		FIELD_PREV_CONTRACTS_BUYBACK: buybackCodes,
		FIELD_PREV_CONTRACTS_SHOP:    shopCodes,
	}, firestore.MergeAll
}

func txDataRemoveManyFromPurchaseQueue[C ICodeAndLocationId](
	remove ...C,
) (map[string]interface{}, firestore.SetOption) {
	m := make(map[int64][]any)
	for _, v := range remove {
		m[v.GetLocationId()] = append(m[v.GetLocationId()], v.GetCode())
	}
	cmd := make(map[string]interface{}, len(m))
	for locationId, codes := range m {
		cmd[fmt.Sprintf("%d", locationId)] = firestore.ArrayRemove(codes...)
	}
	return cmd, firestore.MergeAll
}

func txDataRemoveOneFromPurchaseQueue(
	locationId int64,
	remove string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		fmt.Sprintf("%d", locationId): firestore.ArrayRemove(remove),
	}, firestore.MergeAll
}

func txDataAppendToPurchaseQueue(
	locationId int64,
	append string,
) (map[string]interface{}, firestore.SetOption) {
	return map[string]interface{}{
		fmt.Sprintf("%d", locationId): firestore.ArrayUnion(append),
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
		// FIELD_SHOP_APPRAISAL_REJECTED:     appraisal.Rejected,
		FIELD_SHOP_APPRAISAL_TIME:         appraisal.Time,
		FIELD_SHOP_APPRAISAL_ITEMS:        appraisal.Items,
		FIELD_SHOP_APPRAISAL_VERSION:      appraisal.Version,
		FIELD_SHOP_APPRAISAL_CHARACTER_ID: appraisal.CharacterId,
		FIELD_SHOP_APPRAISAL_LOCATION_ID:  appraisal.LocationId,
		FIELD_SHOP_APPRAISAL_PRICE:        appraisal.Price,
		FIELD_SHOP_APPRAISAL_TAX:          appraisal.Tax,
		FIELD_SHOP_APPRAISAL_TAX_RATE:     appraisal.TaxRate,
	}
}

func txDataSetBuybackAppraisal(
	appraisal BuybackAppraisal,
) map[string]interface{} {
	return map[string]interface{}{
		// FIELD_BUYBACK_APPRAISAL_REJECTED:     appraisal.Rejected,
		FIELD_BUYBACK_APPRAISAL_TIME:         appraisal.Time,
		FIELD_BUYBACK_APPRAISAL_ITEMS:        appraisal.Items,
		FIELD_BUYBACK_APPRAISAL_VERSION:      appraisal.Version,
		FIELD_BUYBACK_APPRAISAL_CHARACTER_ID: appraisal.CharacterId,
		FIELD_BUYBACK_APPRAISAL_SYSTEM_ID:    appraisal.SystemId,
		FIELD_BUYBACK_APPRAISAL_PRICE:        appraisal.Price,
		FIELD_BUYBACK_APPRAISAL_TAX:          appraisal.Tax,
		FIELD_BUYBACK_APPRAISAL_TAX_RATE:     appraisal.TaxRate,
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

func (c *fsClient) readPrevContracts(
	ctx context.Context,
) (
	rep PreviousContracts,
	err error,
) {
	fc, err := c.innerClient()
	if err != nil {
		return rep, err
	}

	ref := prevContractsRef(fc)
	var repPtr *PreviousContracts
	repPtr, err = read[PreviousContracts](ctx, fc, ref)
	if repPtr != nil {
		rep = *repPtr
	}
	return rep, err
}

func (c *fsClient) readPurchaseQueue(
	ctx context.Context,
) (
	rep fsPurchaseQueue,
	err error,
) {
	fc, err := c.innerClient()
	if err != nil {
		return rep, err
	}

	ref := purchaseQueueRef(fc)
	var repPtr *fsPurchaseQueue
	repPtr, err = read[fsPurchaseQueue](ctx, fc, ref)
	if repPtr != nil {
		rep = *repPtr
	}
	return rep, err
}

func (c *fsClient) readUserData(
	ctx context.Context,
	characterId int32,
) (
	rep UserData,
	err error,
) {
	fc, err := c.innerClient()
	if err != nil {
		return rep, err
	}

	ref := userDataRef(characterId, fc)
	var repPtr *UserData
	repPtr, err = read[UserData](ctx, fc, ref)
	if repPtr != nil {
		rep = *repPtr
	}
	return rep, err
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

func (c *fsClient) setPrevContracts(
	ctx context.Context,
	buybackCodes []string,
	shopCodes []string,
) error {
	fc, err := c.innerClient()
	if err != nil {
		return err
	}

	txFunc := func(ctx context.Context, tx *firestore.Transaction) error {
		txpcData, txpcOpts := txDataSetPrevContracts(
			buybackCodes,
			shopCodes,
		)
		if err := txSetPrevContracts(
			fc,
			tx,
			txpcData,
			txpcOpts,
		); err != nil {
			return err
		}
		return nil
	}
	return fc.RunTransaction(ctx, txFunc)
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
		if appraisal.CharacterId != nil {
			// Append the appraisal to user data
			txudData, txudOpts := txDataAppendShopAppraisalToUserData(
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

			// Append the appraisal code to shop queue
			txsqData, txsqOpts := txDataAppendToPurchaseQueue(
				appraisal.LocationId,
				appraisal.Code,
			)
			if err := txSetPurchaseQueue(
				fc,
				tx,
				txsqData,
				txsqOpts,
			); err != nil {
				return err
			}
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

func delShopPurchases[C ICodeAndLocationId](
	c *fsClient,
	ctx context.Context,
	appraisalCodes ...C,
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
	locationId int64,
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
		txsqData, txsqOpts := txDataRemoveOneFromPurchaseQueue(
			locationId,
			appraisalCode,
		)
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
