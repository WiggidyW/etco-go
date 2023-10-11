package purchase

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/chanresult"

	// "github.com/WiggidyW/etco-go/client/appraisal/shop"
	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/client/inventory"
	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	"github.com/WiggidyW/etco-go/client/shopqueue"
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/util"
)

type MakePurchaseResponse struct {
	Status    MakePurchaseStatus
	Appraisal *rdb.ShopAppraisal // ((sometimes)) nil unless Success
}

type MakePurchaseParams struct {
	Items       []appraisal.BasicItem
	LocationId  int64
	CharacterId int32
	Cooldown    time.Duration // time to wait before allowing the purchase
	MaxActive   int           // max number of active purchases allowed
}

type MakePurchaseClient struct {
	appraisalClient     appraisal.MakeShopAppraisalClient
	userDataClient      rdbc.SC_ReadUserDataClient
	writePurchaseClient rdbc.SMAC_WritePurchaseClient
	inventoryClient     inventory.InventoryClient
}

func NewMakePurchaseClient(
	appraisalClient appraisal.MakeShopAppraisalClient,
	userDataClient rdbc.SC_ReadUserDataClient,
	writePurchaseClient rdbc.SMAC_WritePurchaseClient,
	inventoryClient inventory.InventoryClient,
) MakePurchaseClient {
	return MakePurchaseClient{
		appraisalClient,
		userDataClient,
		writePurchaseClient,
		inventoryClient,
	}
}

func (mpc MakePurchaseClient) Fetch(
	ctx context.Context,
	params MakePurchaseParams,
) (*MakePurchaseResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the inventory and shop queue
	chnInvSend, chnInvRecv := chanresult.
		NewChanResult[map[int32]int64](ctx, 0, 0).Split()
	chnSQRepSend, chnSQRepRecv := chanresult.
		NewChanResult[*shopqueue.ShopQueueResponse](ctx, 0, 0).Split()
	go mpc.fetchShopQueueAndInventory(ctx, params, chnSQRepSend, chnInvSend)

	// fetch the appraisal
	chnApprSend, chnApprRecv := chanresult.
		NewChanResult[*rdb.ShopAppraisal](ctx, 0, 0).Split()
	go mpc.fetchAppraisal(ctx, params, chnApprSend)

	// fetch the user data
	userDataPtr, err := mpc.fetchUserData(ctx, params)
	if err != nil {
		return nil, err
	}
	userData := *userDataPtr

	// check if the user is limited by cooldown
	if makeCooldownLimited(userData, params.Cooldown) {
		return &MakePurchaseResponse{
			Status: MPS_CooldownLimit,
			// Appraisal: nil,
		}, nil
	}

	// wait for the shop queue
	sqRep, err := chnSQRepRecv.Recv()
	if err != nil {
		return nil, err
	}
	sqHashSet := sqRep.ShopQueueHashSet()

	// check if the user is limited by max active
	if maxActiveLimited(userData, sqHashSet, params.MaxActive) {
		return &MakePurchaseResponse{
			Status: MPS_MaxActiveLimit,
			// Appraisal: nil,
		}, nil
	}

	// wait for the appraisal and inventory
	appraisal, err := chnApprRecv.Recv()
	if err != nil {
		return nil, err
	}
	inventory, err := chnInvRecv.Recv()
	if err != nil {
		return nil, err
	}

	// check if the purchase is rejected or unavailable
	status := checkReject(appraisal, inventory)
	if status != MPS_Success {
		return &MakePurchaseResponse{
			Status:    status,
			Appraisal: appraisal,
		}, nil
	}

	// write the purchase
	if err := mpc.writePurchase(ctx, appraisal); err != nil {
		return nil, err
	}

	return &MakePurchaseResponse{
		Status:    MPS_Success,
		Appraisal: appraisal,
	}, nil
}

func (mpc MakePurchaseClient) fetchShopQueueAndInventory(
	ctx context.Context,
	params MakePurchaseParams,
	chnSQRepSend chanresult.ChanSendResult[*shopqueue.ShopQueueResponse],
	chnInvSend chanresult.ChanSendResult[map[int32]int64],
) error {
	if inventoryRep, err := mpc.inventoryClient.Fetch(
		ctx,
		inventory.InventoryParams{
			LocationId:          params.LocationId,
			ChnSendShopQueueRep: &chnSQRepSend,
		},
	); err != nil {
		return chnInvSend.SendErr(err)
	} else {
		return chnInvSend.SendOk(*inventoryRep)
	}
}

func (mpc MakePurchaseClient) fetchAppraisal(
	ctx context.Context,
	params MakePurchaseParams,
	chnApprSend chanresult.ChanSendResult[*rdb.ShopAppraisal],
) error {
	if appraisalRep, err := mpc.appraisalClient.Fetch(
		ctx,
		appraisal.MakeShopAppraisalParams{
			Items:       params.Items,
			LocationId:  params.LocationId,
			CharacterId: params.CharacterId,
			IncludeCode: true,
		},
	); err != nil {
		return chnApprSend.SendErr(err)
	} else {
		return chnApprSend.SendOk(appraisalRep)
	}
}

func (mpc MakePurchaseClient) fetchUserData(
	ctx context.Context,
	params MakePurchaseParams,
	// chnSend util.ChanSendResult[rdba.UserData],
) (*rdb.UserData, error) {
	if userDataRep, err := mpc.userDataClient.Fetch(
		ctx,
		rdbc.ReadUserDataParams{CharacterId: params.CharacterId},
	); err != nil {
		return nil, err
	} else {
		userDataVal := userDataRep.Data()
		return &userDataVal, nil
	}
}

// func (mpc MakePurchaseClient) fetchUserData(
// 	ctx context.Context,
// 	params MakePurchaseParams,
// 	chnSend util.ChanSendResult[rdba.UserData],
// ) error {
// 	if userDataRep, err := mpc.UserDataClient.Fetch(
// 		ctx,
// 		rus.ReadUserDataParams{CharacterId: params.CharacterId},
// 	); err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(userDataRep.Data())
// 	}
// }

func (mpc MakePurchaseClient) writePurchase(
	ctx context.Context,
	appraisal *rdb.ShopAppraisal,
) error {
	if timestamp, err := mpc.writePurchaseClient.Fetch(
		ctx,
		rdbc.WritePurchaseParams{Appraisal: *appraisal},
	); err != nil {
		return err
	} else {
		appraisal.Time = *timestamp
		return nil
	}
}

func checkReject(
	appraisal *rdb.ShopAppraisal, // mutates the appraisal
	inventory map[int32]int64,
) MakePurchaseStatus {
	var rejected bool = false
	var unavailable bool = false

	for i := 0; i < len(appraisal.Items); i++ {
		item := &appraisal.Items[i]
		if item.PricePerUnit <= 0 {
			rejected = true
			continue
		}
		quantity, ok := inventory[item.TypeId]
		if !ok || quantity < item.Quantity {
			unavailable = true
			item.Description = fmt.Sprintf(
				"Rejected - %d are available for purchase",
				quantity,
			)
			item.PricePerUnit = 0
		}
	}

	// TODO: check if compiler optimizes this (it's more readable this way)
	if rejected && unavailable {
		return MPS_ItemsRejectedAndUnavailable
	} else if rejected {
		return MPS_ItemsRejected
	} else if unavailable {
		return MPS_ItemsUnavailable
	} else {
		return MPS_Success
	}
}

func makeCooldownLimited(userData rdb.UserData, cooldown time.Duration) bool {
	return time.Now().Before(userData.MadePurchase.Add(cooldown))
}

func maxActiveLimited[HS util.HashSet[string]](
	userData rdb.UserData,
	sqHashSet HS,
	maxActive int,
) bool {
	var active int = 0
	for _, prevPurchase := range userData.ShopAppraisals {
		if sqHashSet.Has(prevPurchase) {
			active++
			if active >= maxActive {
				return true
			}
		}
	}
	return false
}
