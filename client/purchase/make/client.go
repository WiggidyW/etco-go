package make

import (
	"context"
	"fmt"
	"time"

	a "github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	"github.com/WiggidyW/eve-trading-co-go/client/appraisal/shop"
	"github.com/WiggidyW/eve-trading-co-go/client/inventory"
	rdba "github.com/WiggidyW/eve-trading-co-go/client/remotedb/appraisal"
	rus "github.com/WiggidyW/eve-trading-co-go/client/remotedb/appraisal/readuserdata"
	"github.com/WiggidyW/eve-trading-co-go/client/remotedb/purchase/write"
	"github.com/WiggidyW/eve-trading-co-go/client/shopqueue"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type MakePurchaseClient struct {
	AppraisalClient     shop.ShopAppraisalClient
	UserDataClient      rus.SC_ReadUserDataClient
	WritePurchaseClient write.SMAC_WriteShopPurchaseClient
	InventoryClient     inventory.InventoryClient
}

func (mpc MakePurchaseClient) Fetch(
	ctx context.Context,
	params MakePurchaseParams,
) (*MakePurchaseResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the inventory and shop queue
	chnInv := util.NewChanResult[map[int32]int64](ctx)
	chnInvSend, chnInvRecv := chnInv.Split()
	chnSQRep := util.NewChanResult[*shopqueue.ShopQueueResponse](ctx)
	chnSQRepSend, chnSQRepRecv := chnSQRep.Split()
	go mpc.fetchShopQueueAndInventory(ctx, params, chnSQRepSend, chnInvSend)

	// fetch the appraisal
	chnAppr := util.NewChanResult[*a.ShopAppraisal](ctx)
	chnApprSend, chnApprRecv := chnAppr.Split()
	go mpc.fetchAppraisal(ctx, params, chnApprSend)

	// fetch the user data
	userDataPtr, err := mpc.fetchUserData(ctx, params)
	if err != nil {
		return nil, err
	}
	userData := *userDataPtr

	// check if the user is limited by cooldown
	if cooldownLimited(userData, params.Cooldown) {
		return &MakePurchaseResponse{
			Status:    CooldownLimit,
			Appraisal: nil,
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
			Status:    MaxActiveLimit,
			Appraisal: nil,
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
	if status != Success {
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
		Status:    Success,
		Appraisal: appraisal,
	}, nil
}

func (mpc MakePurchaseClient) fetchShopQueueAndInventory(
	ctx context.Context,
	params MakePurchaseParams,
	chnSQRepSend util.ChanSendResult[*shopqueue.ShopQueueResponse],
	chnInvSend util.ChanSendResult[map[int32]int64],
) error {
	if inventoryRep, err := mpc.InventoryClient.Fetch(
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
	chnApprSend util.ChanSendResult[*a.ShopAppraisal],
) error {
	if appraisalRep, err := mpc.AppraisalClient.Fetch(
		ctx,
		shop.ShopAppraisalParams{
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
) (*rdba.UserData, error) {
	if userDataRep, err := mpc.UserDataClient.Fetch(
		ctx,
		rus.ReadUserDataParams{CharacterId: params.CharacterId},
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
	appraisal *a.ShopAppraisal,
) error {
	if timestamp, err := mpc.WritePurchaseClient.Fetch(
		ctx,
		write.WriteShopPurchaseParams{Appraisal: *appraisal},
	); err != nil {
		return err
	} else {
		appraisal.Time = *timestamp
		return nil
	}
}

func cooldownLimited(userData rdba.UserData, cooldown time.Duration) bool {
	return time.Now().Before(userData.MadePurchase.Add(cooldown))
}

func maxActiveLimited[HS util.HashSet[string]](
	userData rdba.UserData,
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

func checkReject(
	appraisal *a.ShopAppraisal, // mutates the appraisal
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
		return ItemsRejectedAndUnavailable
	} else if rejected {
		return ItemsRejected
	} else if unavailable {
		return ItemsUnavailable
	} else {
		return Success
	}
}
