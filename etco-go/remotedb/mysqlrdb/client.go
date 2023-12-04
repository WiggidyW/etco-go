package mysqlrdb

import (
	"context"
	"database/sql"
	"sync"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/remotedb/implrdb"

	_ "github.com/go-sql-driver/mysql"
)

type PreviousContracts = implrdb.PreviousContracts
type HaulAppraisal = implrdb.HaulAppraisal
type HaulItem = implrdb.HaulItem
type ShopAppraisal = implrdb.ShopAppraisal
type ShopItem = implrdb.ShopItem
type BuybackAppraisal = implrdb.BuybackAppraisal
type BuybackParentItem = implrdb.BuybackParentItem
type BuybackChildItem = implrdb.BuybackChildItem
type RawPurchaseQueue = implrdb.RawPurchaseQueue
type UserData = implrdb.UserData
type CodeAndLocationId = implrdb.CodeAndLocationId

type mysqlClient struct {
	_client *sql.DB
	host    string
	mu      *sync.Mutex
}

func NewMySQLClient() *mysqlClient {
	return newMySQLClient(build.RDB_MYSQL_HOST)
}

func newMySQLClient(host string) *mysqlClient {
	return &mysqlClient{
		_client: nil,
		host:    host,
		mu:      new(sync.Mutex),
	}
}

func (c *mysqlClient) innerClient() (*sql.DB, error) {
	if c._client == nil {
		// lock to prevent multiple clients from being created
		c.mu.Lock()
		defer c.mu.Unlock()

		// check again in case another client was created while waiting
		if c._client != nil {
			return c._client, nil
		}

		// create the client
		var err error
		c._client, err = sql.Open("mysql", c.host)
		if err != nil {
			return nil, err
		}

		// initialize the tables
		var tx Transaction
		tx, err = c.beginWriteTx(context.Background())
		if err == nil {
			err = tx.init()
		}
		if err != nil {
			c._client = nil
			return nil, err
		}
	}
	return c._client, nil
}

func (c *mysqlClient) beginWriteTx(
	ctx context.Context,
) (
	tx Transaction,
	err error,
) {
	tx.ctx = ctx
	var client *sql.DB
	client, err = c.innerClient()
	if err != nil {
		return tx, err
	}
	tx.tx, err = client.BeginTx(ctx, nil)
	return tx, err
}

func (c *mysqlClient) beginReadTx(
	ctx context.Context,
) (
	tx Transaction,
	err error,
) {
	tx.ctx = ctx
	var client *sql.DB
	client, err = c.innerClient()
	if err != nil {
		return tx, err
	}
	tx.tx, err = client.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	return tx, err
}

type Transaction struct {
	tx  *sql.Tx
	ctx context.Context
}

func (tx Transaction) exec(
	query string,
	args ...any,
) (sql.Result, error) {
	return tx.tx.ExecContext(tx.ctx, query, args...)
}

func (tx Transaction) query(
	query string,
	args ...any,
) (*sql.Rows, error) {
	return tx.tx.QueryContext(tx.ctx, query, args...)
}

func (tx Transaction) rollback() error {
	return tx.tx.Rollback()
}

func (tx Transaction) commit() error {
	return tx.tx.Commit()
}

func (tx Transaction) init() (err error) {
	_, err = tx.exec(
		`
		CREATE TABLE IF NOT EXISTS b_appraisal (
			b_appraisal_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			rejected BOOLEAN NOT NULL,
			code CHAR(16) NOT NULL,
			time BIGINT NOT NULL,
			version TINYTEXT NOT NULL,
			character_id INT,
			system_id INT NOT NULL,
			price DOUBLE NOT NULL,
			tax DOUBLE,
			tax_rate DOUBLE,
			fee DOUBLE,
			fee_per_m3 DOUBLE
		);
		CREATE TABLE IF NOT EXISTS b_parent_item (
			b_parent_item_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			type_id INT NOT NULL,
			quantity INT NOT NULL,
			price_per_unit DOUBLE NOT NULL,
			description TEXT NOT NULL,
			fee_per_unit DOUBLE
		);
		CREATE TABLE IF NOT EXISTS b_appraisal_b_parent_item (
			b_appraisal_id INT NOT NULL,
			b_parent_item_id INT NOT NULL,
			FOREIGN KEY (b_appraisal_id) REFERENCES b_appraisal(b_appraisal_id),
			FOREIGN KEY (b_parent_item_id) REFERENCES b_parent_item(b_parent_item_id),
			PRIMARY KEY (b_appraisal_id, b_parent_item_id)
		);
		CREATE TABLE IF NOT EXISTS b_child_item (
			b_child_item_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			type_id INT NOT NULL,
			quantity_per_parent DOUBLE NOT NULL,
			price_per_unit DOUBLE NOT NULL,
			description TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS b_parent_item_b_child_item (
			b_parent_item_id INT NOT NULL,
			b_child_item_id INT NOT NULL,
			FOREIGN KEY (b_parent_item_id) REFERENCES b_parent_item(b_parent_item_id),
			FOREIGN KEY (b_child_item_id) REFERENCES b_child_item(b_child_item_id),
			PRIMARY KEY (b_parent_item_id, b_child_item_id)
		);
		CREATE TABLE IF NOT EXISTS s_appraisal (
			s_appraisal_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			rejected BOOLEAN NOT NULL,
			code CHAR(16) NOT NULL,
			time BIGINT NOT NULL,
			version TINYTEXT NOT NULL,
			character_id INT,
			location_id BIGINT NOT NULL,
			price DOUBLE NOT NULL,
			tax DOUBLE,
			tax_rate DOUBLE
		);
		CREATE TABLE IF NOT EXISTS s_item (
			s_item_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			type_id INT NOT NULL,
			quantity INT NOT NULL,
			price_per_unit DOUBLE NOT NULL,
			description TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS s_appraisal_s_item (
			s_appraisal_id INT NOT NULL,
			s_item_id INT NOT NULL,
			FOREIGN KEY (s_appraisal_id) REFERENCES s_appraisal(s_appraisal_id),
			FOREIGN KEY (s_item_id) REFERENCES s_item(s_item_id),
			PRIMARY KEY (s_appraisal_id, s_item_id)
		);
		CREATE TABLE IF NOT EXISTS h_appraisal (
			h_appraisal_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			rejected BOOLEAN NOT NULL,
			code CHAR(16) NOT NULL,
			time BIGINT NOT NULL,
			version TINYTEXT NOT NULL,
			character_id INT,
			start_system_id INT NOT NULL,
			end_system_id INT NOT NULL,
			price DOUBLE NOT NULL,
			tax DOUBLE,
			tax_rate DOUBLE,
			fee_per_m3 DOUBLE,
			collateral_rate DOUBLE,
			reward DOUBLE NOT NULL,
			reward_kind TINYINT UNSIGNED NOT NULL
		);
		CREATE TABLE IF NOT EXISTS h_item (
			h_item_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			type_id INT NOT NULL,
			quantity INT NOT NULL,
			price_per_unit DOUBLE NOT NULL,
			description TEXT NOT NULL,
			fee_per_unit DOUBLE
		);
		CREATE TABLE IF NOT EXISTS h_appraisal_h_item (
			h_appraisal_id INT NOT NULL,
			h_item_id INT NOT NULL,
			FOREIGN KEY (h_appraisal_id) REFERENCES h_appraisal(h_appraisal_id),
			FOREIGN KEY (h_item_id) REFERENCES h_item(h_item_id),
			PRIMARY KEY (h_appraisal_id, h_item_id)
		);
		CREATE TABLE IF NOT EXISTS user_buyback_appraisal (
			character_id INT NOT NULL PRIMARY KEY,
			code CHAR(16) NOT NULL,
			FOREIGN KEY (code) REFERENCES b_appraisal(code)
		);
		CREATE TABLE IF NOT EXISTS user_shop_appraisal (
			character_id INT NOT NULL PRIMARY KEY,
			code CHAR(16) NOT NULL,
			FOREIGN KEY (code) REFERENCES s_appraisal(code)
		);
		CREATE_TABLE_IF_NOT_EXISTS user_haul_appraisal (
			character_id INT NOT NULL PRIMARY KEY,
			code CHAR(16) NOT NULL,
			FOREIGN KEY (code) REFERENCES h_appraisal(code)
		);
		CREATE TABLE IF NOT EXISTS user_made_purchase (
			character_id INT NOT NULL PRIMARY KEY,
			time BIGINT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS user_cancelled_purchase (
			character_id INT NOT NULL PRIMARY KEY,
			time BIGINT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS purchase_queue (
			code CHAR(16) NOT NULL PRIMARY KEY,
			location_id BIGINT NOT NULL,
			FOREIGN KEY (code) REFERENCES s_appraisal(code)
		);
		CREATE TABLE IF NOT EXISTS prev_buyback_contract (
			code CHAR(16) NOT NULL PRIMARY KEY
		);
		CREATE TABLE IF NOT EXISTS prev_shop_contract (
			code CHAR(16) NOT NULL PRIMARY KEY
		);
		CREATE TABLE IF NOT EXISTS prev_haul_contract (
			code CHAR(16) NOT NULL PRIMARY KEY
		);
		`,
	)
	return err
}

func (tx Transaction) selectBuybackAppraisal(
	code string,
) (
	bAppraisalId int64,
	bAppraisal *BuybackAppraisal,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(
		`
		SELECT
			b_appraisal_id,
			rejected,
			time,
			version,
			character_id,
			system_id,
			price,
			tax,
			tax_rate,
			fee,
			fee_per_m3
		FROM b_appraisal
		WHERE code = ?;
		`,
		code,
	)
	if err != nil || !rows.Next() {
		return 0, nil, err
	}
	var (
		rejected    bool
		timestamp   int64
		version     string
		characterId *int32
		systemId    int32
		price       float64
		tax         float64
		taxRate     float64
		fee         float64
		feePerM3    float64
	)
	err = rows.Scan(
		&bAppraisalId,
		&rejected,
		&timestamp,
		&version,
		characterId,
		&systemId,
		&price,
		&tax,
		&taxRate,
		&fee,
		&feePerM3,
	)
	if err != nil {
		return 0, nil, err
	}
	bAppraisal = &BuybackAppraisal{
		Rejected:    rejected,
		Code:        code,
		Time:        time.Unix(timestamp, 0),
		Version:     version,
		CharacterId: characterId,
		SystemId:    systemId,
		Price:       price,
		Tax:         tax,
		TaxRate:     taxRate,
		Fee:         fee,
		FeePerM3:    feePerM3,
	}
	return bAppraisalId, bAppraisal, nil
}

func (tx Transaction) selectBuybackParentItems(
	bAppraisalId int64,
) (
	bParentItemIds []int64,
	bParentItems []BuybackParentItem,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(
		`
		SELECT
			b_parent_item_id,
			type_id,
			quantity,
			price_per_unit,
			description,
			fee_per_unit
		FROM b_parent_item
		INNER JOIN b_appraisal_b_parent_item
		ON b_parent_item.b_parent_item_id = b_appraisal_b_parent_item.b_parent_item_id
		WHERE b_appraisal_id = ?;
		`,
		bAppraisalId,
	)
	if err != nil {
		return nil, nil, err
	}
	bParentItemIds = make([]int64, 0)
	bParentItems = make([]BuybackParentItem, 0)
	for rows.Next() {
		var (
			bParentItemId int64
			typeId        int32
			quantity      int64
			pricePerUnit  float64
			description   string
			feePerUnit    float64
		)
		err = rows.Scan(
			&bParentItemId,
			&typeId,
			&quantity,
			&pricePerUnit,
			&description,
			&feePerUnit,
		)
		if err != nil {
			return nil, nil, err
		}
		bParentItemIds = append(bParentItemIds, bParentItemId)
		bParentItems = append(bParentItems, BuybackParentItem{
			TypeId:       typeId,
			Quantity:     quantity,
			PricePerUnit: pricePerUnit,
			Description:  description,
			FeePerUnit:   feePerUnit,
		})
	}
	return bParentItemIds, bParentItems, nil
}

func (tx Transaction) selectBuybackChildItems(
	bParentItemId int64,
) (
	bChildItems []BuybackChildItem,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(
		`
		SELECT
			type_id,
			quantity_per_parent,
			price_per_unit,
			description
		FROM b_child_item
		INNER JOIN b_parent_item_b_child_item
		ON b_child_item.b_child_item_id = b_parent_item_b_child_item.b_child_item_id
		WHERE b_parent_item_id = ?;
		`,
		bParentItemId,
	)
	if err != nil {
		return nil, err
	}
	bChildItems = make([]BuybackChildItem, 0)
	for rows.Next() {
		var (
			typeId            int32
			quantityPerParent float64
			pricePerUnit      float64
			description       string
		)
		err = rows.Scan(
			&typeId,
			&quantityPerParent,
			&pricePerUnit,
			&description,
		)
		if err != nil {
			return nil, err
		}
		bChildItems = append(bChildItems, BuybackChildItem{
			TypeId:            typeId,
			QuantityPerParent: quantityPerParent,
			PricePerUnit:      pricePerUnit,
			Description:       description,
		})
	}
	return bChildItems, nil
}

func (tx Transaction) selectShopAppraisal(
	code string,
) (
	sAppraisalId int64,
	sAppraisal *ShopAppraisal,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(
		`
		SELECT
			s_appraisal_id,
			rejected,
			time,
			version,
			character_id,
			location_id,
			price,
			tax,
			tax_rate
		FROM s_appraisal
		WHERE code = ?;
		`,
		code,
	)
	if err != nil || !rows.Next() {
		return 0, nil, err
	}
	var (
		rejected    bool
		timestamp   int64
		version     string
		characterId *int32
		locationId  int64
		price       float64
		tax         float64
		taxRate     float64
	)
	err = rows.Scan(
		&sAppraisalId,
		&rejected,
		&timestamp,
		&version,
		characterId,
		&locationId,
		&price,
		&tax,
		&taxRate,
	)
	if err != nil {
		return 0, nil, err
	}
	sAppraisal = &ShopAppraisal{
		Rejected:    rejected,
		Code:        code,
		Time:        time.Unix(timestamp, 0),
		Version:     version,
		CharacterId: characterId,
		LocationId:  locationId,
		Price:       price,
		Tax:         tax,
		TaxRate:     taxRate,
	}
	return sAppraisalId, sAppraisal, nil
}

func (tx Transaction) selectShopItems(
	sAppraisalId int64,
) (
	sItems []ShopItem,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(
		`
		SELECT
			s_item_id,
			type_id,
			quantity,
			price_per_unit,
			description
		FROM s_item
		INNER JOIN s_appraisal_s_item
		ON s_item.s_item_id = s_appraisal_s_item.s_item_id
		WHERE s_appraisal_id = ?;
		`,
		sAppraisalId,
	)
	if err != nil {
		return nil, err
	}
	sItems = make([]ShopItem, 0)
	for rows.Next() {
		var (
			sItemId      int64
			typeId       int32
			quantity     int64
			pricePerUnit float64
			description  string
		)
		err = rows.Scan(
			&sItemId,
			&typeId,
			&quantity,
			&pricePerUnit,
			&description,
		)
		if err != nil {
			return nil, err
		}
		sItems = append(sItems, ShopItem{
			TypeId:       typeId,
			Quantity:     quantity,
			PricePerUnit: pricePerUnit,
			Description:  description,
		})
	}
	return sItems, nil
}

func (tx Transaction) selectHaulAppraisal(
	code string,
) (
	hAppraisalId int64,
	hAppraisal *HaulAppraisal,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(
		`
		SELECT
			h_appraisal_id,
			rejected,
			time,
			version,
			character_id,
			start_system_id,
			end_system_id,
			price,
			tax,
			tax_rate,
			fee_per_m3,
			collateral_rate,
			reward,
			reward_kind
		FROM h_appraisal
		WHERE code = ?;
		`,
		code,
	)
	if err != nil || !rows.Next() {
		return 0, nil, err
	}
	var (
		rejected       bool
		timestamp      int64
		version        string
		characterId    *int32
		startSystemId  int32
		endSystemId    int32
		price          float64
		tax            float64
		taxRate        float64
		feePerM3       float64
		collateralRate float64
		reward         float64
		rewardKind     uint8
	)
	err = rows.Scan(
		&hAppraisalId,
		&rejected,
		&timestamp,
		&version,
		characterId,
		&startSystemId,
		&endSystemId,
		&price,
		&tax,
		&taxRate,
		&feePerM3,
		&collateralRate,
		&reward,
		&rewardKind,
	)
	if err != nil {
		return 0, nil, err
	}
	hAppraisal = &HaulAppraisal{
		Rejected:       rejected,
		Code:           code,
		Time:           time.Unix(timestamp, 0),
		Version:        version,
		CharacterId:    characterId,
		StartSystemId:  startSystemId,
		EndSystemId:    endSystemId,
		Price:          price,
		Tax:            tax,
		TaxRate:        taxRate,
		FeePerM3:       feePerM3,
		CollateralRate: collateralRate,
		Reward:         reward,
		RewardKind:     rewardKind,
	}
	return hAppraisalId, hAppraisal, nil
}

func (tx Transaction) selectHaulItems(
	hAppraisalId int64,
) (
	hItems []HaulItem,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(
		`
		SELECT
			h_item_id,
			type_id,
			quantity,
			price_per_unit,
			description,
			fee_per_unit
		FROM h_item
		INNER JOIN h_appraisal_h_item
		ON h_item.h_item_id = h_appraisal_h_item.h_item_id
		WHERE h_appraisal_id = ?;
		`,
		hAppraisalId,
	)
	if err != nil {
		return nil, err
	}
	hItems = make([]HaulItem, 0)
	for rows.Next() {
		var (
			hItemId      int64
			typeId       int32
			quantity     int64
			pricePerUnit float64
			description  string
			feePerUnit   float64
		)
		err = rows.Scan(
			&hItemId,
			&typeId,
			&quantity,
			&pricePerUnit,
			&description,
			&feePerUnit,
		)
		if err != nil {
			return nil, err
		}
		hItems = append(hItems, HaulItem{
			TypeId:       typeId,
			Quantity:     quantity,
			PricePerUnit: pricePerUnit,
			Description:  description,
			FeePerUnit:   feePerUnit,
		})
	}
	return hItems, nil
}

func (tx Transaction) selectPurchaseQueue() (
	purchaseQueue RawPurchaseQueue,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(
		`
		SELECT code, location_id
		FROM purchase_queue;
		`,
	)
	if err != nil {
		return nil, err
	}
	purchaseQueue = make(RawPurchaseQueue)
	for rows.Next() {
		var (
			code       string
			locationId int64
		)
		err = rows.Scan(&code, &locationId)
		if err != nil {
			return nil, err
		}
		purchaseQueue[locationId] = append(purchaseQueue[locationId], code)
	}
	return purchaseQueue, nil
}

const (
	SELECT_USER_B_APPRAISALS string = `
		SELECT code
		FROM user_buyback_appraisal
		WHERE character_id = ?;
	`
	SELECT_USER_S_APPRAISALS string = `
		SELECT code
		FROM user_shop_appraisal
		WHERE character_id = ?;
	`
	SELECT_USER_H_APPRAISALS string = `
		SELECT code
		FROM user_haul_appraisal
		WHERE character_id = ?;
	`
)

func (tx Transaction) selectUserAppraisals(
	characterId int32,
	query string,
) (
	codes []string,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(query, characterId)
	if err != nil {
		return nil, err
	}
	codes = make([]string, 0)
	for rows.Next() {
		var code string
		err = rows.Scan(&code)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}

func (tx Transaction) selectUserBuybackAppraisals(characterId int32) (
	codes []string,
	err error,
) {
	return tx.selectUserAppraisals(characterId, SELECT_USER_B_APPRAISALS)
}

func (tx Transaction) selectUserShopAppraisals(characterId int32) (
	codes []string,
	err error,
) {
	return tx.selectUserAppraisals(characterId, SELECT_USER_S_APPRAISALS)
}

func (tx Transaction) selectUserHaulAppraisals(characterId int32) (
	codes []string,
	err error,
) {
	return tx.selectUserAppraisals(characterId, SELECT_USER_H_APPRAISALS)
}

const (
	SELECT_USER_MADE_PURCHASE string = `
		SELECT time
		FROM user_made_purchase
		WHERE character_id = ?;
	`
	SELECT_USER_CANCELLED_PURCHASE string = `
		SELECT time
		FROM user_cancelled_purchase
		WHERE character_id = ?;
	`
)

func (tx Transaction) selectUserTimestamp(
	characterId int32,
	query string,
) (
	timestamp *time.Time,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(query, characterId)
	if err != nil || !rows.Next() {
		return nil, err
	}
	var timestampInt *int64
	err = rows.Scan(timestampInt)
	if err != nil {
		return nil, err
	}
	if timestampInt != nil {
		timestampVal := time.Unix(*timestampInt, 0)
		timestamp = &timestampVal
	}
	return timestamp, nil
}

func (tx Transaction) selectUserMadePurchase(characterId int32) (
	timestamp *time.Time,
	err error,
) {
	return tx.selectUserTimestamp(characterId, SELECT_USER_MADE_PURCHASE)
}

func (tx Transaction) selectUserCancelledPurchase(characterId int32) (
	timestamp *time.Time,
	err error,
) {
	return tx.selectUserTimestamp(characterId, SELECT_USER_CANCELLED_PURCHASE)
}

const (
	SELECT_PREV_B_CONTRACTS string = `
		SELECT code
		FROM prev_buyback_contract;
	`
	SELECT_PREV_S_CONTRACTS string = `
		SELECT code
		FROM prev_shop_contract;
	`
	SELECT_PREV_H_CONTRACTS string = `
		SELECT code
		FROM prev_haul_contract;
	`
)

func (tx Transaction) selectPrevContracts(query string) (
	codes []string,
	err error,
) {
	var rows *sql.Rows
	rows, err = tx.query(query)
	if err != nil {
		return nil, err
	}
	codes = make([]string, 0)
	for rows.Next() {
		var code string
		err = rows.Scan(&code)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}

func (tx Transaction) selectPrevBuybackContracts() (
	codes []string,
	err error,
) {
	return tx.selectPrevContracts(SELECT_PREV_B_CONTRACTS)
}

func (tx Transaction) selectPrevShopContracts() (
	codes []string,
	err error,
) {
	return tx.selectPrevContracts(SELECT_PREV_S_CONTRACTS)
}

func (tx Transaction) selectPrevHaulContracts() (
	codes []string,
	err error,
) {
	return tx.selectPrevContracts(SELECT_PREV_H_CONTRACTS)
}

func (tx Transaction) insertBuybackAppraisal(
	appraisal BuybackAppraisal,
) (
	bAppraisalId int64,
	err error,
) {
	var res sql.Result
	res, err = tx.exec(
		`
		INSERT INTO b_appraisal (
			rejected,
			code,
			time,
			version,
			character_id,
			system_id,
			price,
			tax,
			tax_rate,
			fee,
			fee_per_m3
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
		`,
		appraisal.Rejected,
		appraisal.Code,
		appraisal.Time.Unix(),
		appraisal.Version,
		appraisal.CharacterId,
		appraisal.SystemId,
		appraisal.Price,
		appraisal.Tax,
		appraisal.TaxRate,
		appraisal.Fee,
		appraisal.FeePerM3,
	)
	if err != nil || len(appraisal.Items) < 1 {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

func (tx Transaction) insertBuybackParentItem(
	item BuybackParentItem,
	bAppraisalId int64,
) (
	bParentItemId int64,
	err error,
) {
	var res sql.Result
	res, err = tx.exec(
		`
		INSERT INTO b_parent_item (
			type_id,
			quantity,
			price_per_unit,
			description,
			fee_per_unit
		) VALUES (?, ?, ?, ?, ?);
		`,
		item.TypeId,
		item.Quantity,
		item.PricePerUnit,
		item.Description,
		item.FeePerUnit,
	)
	if err == nil {
		bParentItemId, err = res.LastInsertId()
	}
	if err == nil {
		err = tx.insertBuybackParentItemId(bAppraisalId, bParentItemId)
	}
	return bParentItemId, err
}

func (tx Transaction) insertBuybackChildItem(
	item BuybackChildItem,
	bParentItemId int64,
) (err error) {
	var res sql.Result
	var bChildItemId int64
	res, err = tx.exec(
		`
		INSERT INTO b_child_item (
			type_id,
			quantity_per_parent,
			price_per_unit,
			description
		) VALUES (?, ?, ?, ?);
		`,
		item.TypeId,
		item.QuantityPerParent,
		item.PricePerUnit,
		item.Description,
	)
	if err == nil {
		bChildItemId, err = res.LastInsertId()
	}
	if err == nil {
		err = tx.insertBuybackChildItemId(bParentItemId, bChildItemId)
	}
	return err
}

func (tx Transaction) insertShopAppraisal(
	appraisal ShopAppraisal,
) (
	sAppraisalId int64,
	err error,
) {
	var res sql.Result
	res, err = tx.exec(
		`
		INSERT INTO s_appraisal (
			rejected,
			code,
			time,
			version,
			character_id,
			location_id,
			price,
			tax,
			tax_rate
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
		`,
		appraisal.Rejected,
		appraisal.Code,
		appraisal.Time.Unix(),
		appraisal.Version,
		appraisal.CharacterId,
		appraisal.LocationId,
		appraisal.Price,
		appraisal.Tax,
		appraisal.TaxRate,
	)
	if err != nil || len(appraisal.Items) < 1 {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

func (tx Transaction) insertShopItem(
	item ShopItem,
	sAppraisalId int64,
) (err error) {
	var res sql.Result
	var sItemId int64
	res, err = tx.exec(
		`
		INSERT INTO s_item (
			type_id,
			quantity,
			price_per_unit,
			description
		) VALUES (?, ?, ?, ?);
		`,
		item.TypeId,
		item.Quantity,
		item.PricePerUnit,
		item.Description,
	)
	if err == nil {
		sItemId, err = res.LastInsertId()
	}
	if err == nil {
		err = tx.insertShopItemId(sAppraisalId, sItemId)
	}
	return err
}

func (tx Transaction) insertHaulAppraisal(
	appraisal HaulAppraisal,
) (
	hAppraisalId int64,
	err error,
) {
	var res sql.Result
	res, err = tx.exec(
		`
		INSERT INTO h_appraisal (
			rejected,
			code,
			time,
			version,
			character_id,
			start_system_id,
			end_system_id,
			price,
			tax,
			tax_rate,
			fee_per_m3,
			collateral_rate,
			reward,
			reward_kind
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
		`,
		appraisal.Rejected,
		appraisal.Code,
		appraisal.Time.Unix(),
		appraisal.Version,
		appraisal.CharacterId,
		appraisal.StartSystemId,
		appraisal.EndSystemId,
		appraisal.Price,
		appraisal.Tax,
		appraisal.TaxRate,
		appraisal.FeePerM3,
		appraisal.Reward,
		appraisal.RewardKind,
	)
	if err != nil || len(appraisal.Items) < 1 {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

func (tx Transaction) insertHaulItem(
	item HaulItem,
	hAppraisalId int64,
) (err error) {
	var res sql.Result
	var hItemId int64
	res, err = tx.exec(
		`
		INSERT INTO h_item (
			type_id,
			quantity,
			price_per_unit,
			description,
			fee_per_unit
		) VALUES (?, ?, ?, ?, ?);
		`,
		item.TypeId,
		item.Quantity,
		item.PricePerUnit,
		item.Description,
		item.FeePerUnit,
	)
	if err == nil {
		hItemId, err = res.LastInsertId()
	}
	if err == nil {
		err = tx.insertHaulItemId(hAppraisalId, hItemId)
	}
	return err
}

const (
	INSERT_B_PARENT_ITEM_ID string = `
		INSERT INTO b_appraisal_b_parent_item (
			b_appraisal_id,
			b_parent_item_id
		) VALUES (?, ?);
	`
	INSERT_B_CHILD_ITEM_ID string = `
		INSERT INTO b_parent_item_b_child_item (
			b_parent_item_id,
			b_child_item_id
		) VALUES (?, ?);
	`
	INSERT_S_ITEM_ID string = `
		INSERT INTO s_appraisal_s_item (
			s_appraisal_id,
			s_item_id
		) VALUES (?, ?);
	`
	INSERT_H_ITEM_ID string = `
		INSERT INTO h_appraisal_h_item (
			h_appraisal_id,
			h_item_id
		) VALUES (?, ?);
	`
)

func (tx Transaction) insertItemId(
	parentId int64,
	itemId int64,
	query string,
) (err error) {
	_, err = tx.exec(query, parentId, itemId)
	return err
}

func (tx Transaction) insertBuybackParentItemId(
	bAppraisalId int64,
	bParentItemId int64,
) (err error) {
	return tx.insertItemId(bAppraisalId, bParentItemId, INSERT_B_PARENT_ITEM_ID)
}

func (tx Transaction) insertBuybackChildItemId(
	bParentItemId int64,
	bChildItemId int64,
) (err error) {
	return tx.insertItemId(bParentItemId, bChildItemId, INSERT_B_CHILD_ITEM_ID)
}

func (tx Transaction) insertShopItemId(
	sAppraisalId int64,
	sItemId int64,
) (err error) {
	return tx.insertItemId(sAppraisalId, sItemId, INSERT_S_ITEM_ID)
}

func (tx Transaction) insertHaulItemId(
	hAppraisalId int64,
	hItemId int64,
) (err error) {
	return tx.insertItemId(hAppraisalId, hItemId, INSERT_H_ITEM_ID)
}

func (tx Transaction) delShopPurchases(
	appraisalCodes ...CodeAndLocationId,
) (err error) {
	for _, codeAndLocationId := range appraisalCodes {
		_, err = tx.exec(
			`
			DELETE FROM purchase_queue
			WHERE code = ? AND location_id = ?;
			`,
			codeAndLocationId.Code,
			codeAndLocationId.LocationId,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	INSERT_USER_B_APPRAISAL string = `
		INSERT INTO user_buyback_appraisal (
			code,
			character_id
		) VALUES (?, ?);
	`
	INSERT_USER_S_APPRAISAL string = `
		INSERT INTO user_shop_appraisal (
			code,
			character_id
		) VALUES (?, ?);
	`
	INSERT_USER_H_APPRAISAL string = `
		INSERT INTO user_haul_appraisal (
			code,
			character_id
		) VALUES (?, ?);
	`
)

func (tx Transaction) insertUserAppraisal(
	code string,
	characterId int32,
	query string,
) (err error) {
	_, err = tx.exec(query, code, characterId)
	return err
}

func (tx Transaction) insertUserBuybackAppraisal(
	code string,
	characterId int32,
) (err error) {
	return tx.insertUserAppraisal(code, characterId, INSERT_USER_B_APPRAISAL)
}

func (tx Transaction) insertUserShopAppraisal(
	code string,
	characterId int32,
) (err error) {
	return tx.insertUserAppraisal(code, characterId, INSERT_USER_S_APPRAISAL)
}

func (tx Transaction) insertUserHaulAppraisal(
	code string,
	characterId int32,
) (err error) {
	return tx.insertUserAppraisal(code, characterId, INSERT_USER_H_APPRAISAL)
}

const (
	INSERT_USER_MADE_PURCHASE string = `
		INSERT INTO user_made_purchase (
			character_id,
			time
		) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			time = VALUES(time);
	`
	INSERT_USER_CANCELLED_PURCHASE string = `
		INSERT INTO user_cancelled_purchase (
			character_id,
			time
		) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			time = VALUES(time);
	`
)

func (tx Transaction) insertUserTimestamp(
	characterId int32,
	timestamp int64,
	query string,
) (err error) {
	_, err = tx.exec(query, characterId, timestamp)
	return err
}

func (tx Transaction) insertUserCancelledPurchase(
	characterId int32,
	timestamp int64,
) (err error) {
	return tx.insertUserTimestamp(
		characterId,
		timestamp,
		INSERT_USER_CANCELLED_PURCHASE,
	)
}

func (tx Transaction) insertUserMadePurchase(
	characterId int32,
	timestamp int64,
) (err error) {
	return tx.insertUserTimestamp(
		characterId,
		timestamp,
		INSERT_USER_MADE_PURCHASE,
	)
}

func (tx Transaction) insertPurchase(
	code string,
	locationId int64,
) (err error) {
	_, err = tx.exec(
		`
		INSERT INTO purchase_queue (
			code,
			location_id
		) VALUES (?, ?);
		`,
		code,
		locationId,
	)
	return err
}

const (
	TRUNCATE_PREV_B_CONTRACTS string = `TRUNCATE TABLE prev_buyback_contract;`
	TRUNCATE_PREV_S_CONTRACTS string = `TRUNCATE TABLE prev_shop_contract;`
	TRUNCATE_PREV_H_CONTRACTS string = `TRUNCATE TABLE prev_haul_contract;`

	INSERT_PREV_B_CONTRACT string = `
		INSERT INTO prev_buyback_contract (
			code
		) VALUES (?);
	`
	INSERT_PREV_S_CONTRACT string = `
		INSERT INTO prev_shop_contract (
			code
		) VALUES (?);
	`
	INSERT_PREV_H_CONTRACT string = `
		INSERT INTO prev_haul_contract (
			code
		) VALUES (?);
	`
)

func (tx Transaction) insertPrevContracts(
	codes []string,
	query string,
	trunc string,
) (err error) {
	_, err = tx.exec(trunc)
	if err != nil {
		return err
	}
	for _, code := range codes {
		_, err = tx.exec(query, code)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tx Transaction) insertPrevBuybackContracts(codes []string) (err error) {
	return tx.insertPrevContracts(
		codes,
		INSERT_PREV_B_CONTRACT,
		TRUNCATE_PREV_B_CONTRACTS,
	)
}

func (tx Transaction) insertPrevShopContracts(codes []string) (err error) {
	return tx.insertPrevContracts(
		codes,
		INSERT_PREV_S_CONTRACT,
		TRUNCATE_PREV_S_CONTRACTS,
	)
}

func (tx Transaction) insertPrevHaulContracts(codes []string) (err error) {
	return tx.insertPrevContracts(
		codes,
		INSERT_PREV_H_CONTRACT,
		TRUNCATE_PREV_H_CONTRACTS,
	)
}

func (c *mysqlClient) ReadBuybackAppraisal(
	ctx context.Context,
	appraisalCode string,
) (
	appraisal *BuybackAppraisal,
	err error,
) {
	var tx Transaction
	tx, err = c.beginReadTx(ctx)
	if err != nil {
		return nil, err
	} else {
		defer tx.rollback()
	}
	var bAppraisalId int64
	bAppraisalId, appraisal, err = tx.selectBuybackAppraisal(appraisalCode)
	if err != nil || appraisal == nil {
		return nil, err
	}
	var bParentItemIds []int64
	bParentItemIds, appraisal.Items, err = tx.selectBuybackParentItems(
		bAppraisalId,
	)
	if err != nil || len(appraisal.Items) < 1 {
		return nil, err
	}
	for i, bParentItemId := range bParentItemIds {
		appraisal.Items[i].Children, err = tx.selectBuybackChildItems(
			bParentItemId,
		)
		if err != nil {
			return nil, err
		}
	}
	return appraisal, nil
}

func (c *mysqlClient) ReadShopAppraisal(
	ctx context.Context,
	appraisalCode string,
) (
	appraisal *ShopAppraisal,
	err error,
) {
	var tx Transaction
	tx, err = c.beginReadTx(ctx)
	if err != nil {
		return nil, err
	} else {
		defer tx.rollback()
	}
	var sAppraisalId int64
	sAppraisalId, appraisal, err = tx.selectShopAppraisal(appraisalCode)
	if err != nil || appraisal == nil {
		return nil, err
	}
	appraisal.Items, err = tx.selectShopItems(sAppraisalId)
	return appraisal, err
}

func (c *mysqlClient) ReadHaulAppraisal(
	ctx context.Context,
	appraisalCode string,
) (
	appraisal *HaulAppraisal,
	err error,
) {
	var tx Transaction
	tx, err = c.beginReadTx(ctx)
	if err != nil {
		return nil, err
	} else {
		defer tx.rollback()
	}
	var hAppraisalId int64
	hAppraisalId, appraisal, err = tx.selectHaulAppraisal(appraisalCode)
	if err != nil || appraisal == nil {
		return nil, err
	}
	appraisal.Items, err = tx.selectHaulItems(hAppraisalId)
	return appraisal, err
}

func (c *mysqlClient) ReadPurchaseQueue(
	ctx context.Context,
) (
	queue RawPurchaseQueue,
	err error,
) {
	var tx Transaction
	tx, err = c.beginReadTx(ctx)
	if err != nil {
		return nil, err
	} else {
		defer tx.rollback()
	}
	return tx.selectPurchaseQueue()
}

func (c *mysqlClient) ReadPrevContracts(
	ctx context.Context,
) (
	contracts PreviousContracts,
	err error,
) {
	var tx Transaction
	tx, err = c.beginReadTx(ctx)
	if err != nil {
		return contracts, err
	} else {
		defer tx.rollback()
	}
	contracts.Buyback, err = tx.selectPrevBuybackContracts()
	if err != nil {
		return contracts, err
	}
	contracts.Shop, err = tx.selectPrevShopContracts()
	if err != nil {
		return contracts, err
	}
	contracts.Haul, err = tx.selectPrevHaulContracts()
	if err != nil {
		return contracts, err
	}
	return contracts, nil
}

func (c *mysqlClient) ReadUserData(
	ctx context.Context,
	characterId int32,
) (
	userData UserData,
	err error,
) {
	var tx Transaction
	tx, err = c.beginReadTx(ctx)
	if err != nil {
		return userData, err
	} else {
		defer tx.rollback()
	}
	userData.BuybackAppraisals, err = tx.selectUserBuybackAppraisals(
		characterId,
	)
	if err != nil {
		return userData, err
	}
	userData.ShopAppraisals, err = tx.selectUserShopAppraisals(characterId)
	if err != nil {
		return userData, err
	}
	userData.HaulAppraisals, err = tx.selectUserHaulAppraisals(characterId)
	if err != nil {
		return userData, err
	}
	userData.MadePurchase, err = tx.selectUserMadePurchase(characterId)
	if err != nil {
		return userData, err
	}
	userData.CancelledPurchase, err = tx.selectUserCancelledPurchase(
		characterId,
	)
	if err != nil {
		return userData, err
	}
	return userData, nil
}

func (c *mysqlClient) SetPrevContracts(
	ctx context.Context,
	buybackCodes []string,
	shopCodes []string,
	haulCodes []string,
) (err error) {
	var tx Transaction
	tx, err = c.beginWriteTx(ctx)
	if err != nil {
		return err
	}
	err = tx.insertPrevBuybackContracts(buybackCodes)
	if err != nil {
		tx.rollback()
		return err
	}
	err = tx.insertPrevShopContracts(shopCodes)
	if err != nil {
		tx.rollback()
		return err
	}
	err = tx.insertPrevHaulContracts(haulCodes)
	if err != nil {
		tx.rollback()
		return err
	}
	return tx.commit()
}

func (c *mysqlClient) DelShopPurchases(
	ctx context.Context,
	appraisalCodes ...CodeAndLocationId,
) (err error) {
	var tx Transaction
	tx, err = c.beginWriteTx(ctx)
	if err != nil {
		return err
	}
	err = tx.delShopPurchases(appraisalCodes...)
	if err != nil {
		tx.rollback()
	} else {
		err = tx.commit()
	}
	return err
}

func (c *mysqlClient) CancelShopPurchase(
	ctx context.Context,
	characterId int32,
	appraisalCode string,
	locationId int64,
) (err error) {
	var tx Transaction
	tx, err = c.beginWriteTx(ctx)
	if err != nil {
		return err
	}
	err = tx.insertUserCancelledPurchase(characterId, time.Now().Unix())
	if err != nil {
		tx.rollback()
		return err
	}
	err = tx.delShopPurchases(CodeAndLocationId{
		Code:       appraisalCode,
		LocationId: locationId,
	})
	if err != nil {
		tx.rollback()
		return err
	}
	return tx.commit()
}

func (c *mysqlClient) SaveBuybackAppraisal(
	ctx context.Context,
	appraisal BuybackAppraisal,
) (err error) {
	var tx Transaction
	var bAppraisalId, bParentItemId int64
	tx, err = c.beginWriteTx(ctx)
	if err != nil {
		return err
	}
	bAppraisalId, err = tx.insertBuybackAppraisal(appraisal)
	if err != nil {
		tx.rollback()
		return err
	}
	for _, item := range appraisal.Items {
		bParentItemId, err = tx.insertBuybackParentItem(
			item,
			bAppraisalId,
		)
		if err != nil {
			tx.rollback()
			return err
		}
		for _, childItem := range item.Children {
			err = tx.insertBuybackChildItem(childItem, bParentItemId)
			if err != nil {
				tx.rollback()
				return err
			}
		}
	}
	if appraisal.CharacterId != nil {
		characterId := *appraisal.CharacterId
		err = tx.insertUserBuybackAppraisal(appraisal.Code, characterId)
		if err != nil {
			tx.rollback()
			return err
		}
	}
	return tx.commit()
}

func (c *mysqlClient) SaveShopAppraisal(
	ctx context.Context,
	appraisal ShopAppraisal,
) (err error) {
	var tx Transaction
	var sAppraisalId int64
	tx, err = c.beginWriteTx(ctx)
	if err != nil {
		return err
	}
	sAppraisalId, err = tx.insertShopAppraisal(appraisal)
	if err != nil {
		tx.rollback()
		return err
	}
	for _, item := range appraisal.Items {
		err = tx.insertShopItem(item, sAppraisalId)
		if err != nil {
			tx.rollback()
			return err
		}
	}
	if appraisal.CharacterId != nil {
		characterId := *appraisal.CharacterId
		err = tx.insertPurchase(appraisal.Code, appraisal.LocationId)
		if err != nil {
			tx.rollback()
			return err
		}
		err = tx.insertUserMadePurchase(characterId, time.Now().Unix())
		if err != nil {
			tx.rollback()
			return err
		}
		err = tx.insertUserShopAppraisal(appraisal.Code, characterId)
		if err != nil {
			tx.rollback()
			return err
		}
	}
	return tx.commit()
}

func (c *mysqlClient) SaveHaulAppraisal(
	ctx context.Context,
	appraisal HaulAppraisal,
) (err error) {
	var tx Transaction
	var hAppraisalId int64
	tx, err = c.beginWriteTx(ctx)
	if err != nil {
		return err
	}
	hAppraisalId, err = tx.insertHaulAppraisal(appraisal)
	if err != nil {
		tx.rollback()
		return err
	}
	for _, item := range appraisal.Items {
		err = tx.insertHaulItem(item, hAppraisalId)
		if err != nil {
			tx.rollback()
			return err
		}
	}
	if appraisal.CharacterId != nil {
		characterId := *appraisal.CharacterId
		err = tx.insertUserHaulAppraisal(appraisal.Code, characterId)
		if err != nil {
			tx.rollback()
			return err
		}
	}
	return tx.commit()
}
