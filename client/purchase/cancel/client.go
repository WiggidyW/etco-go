package cancel

import (
	"context"
	"time"

	rdba "github.com/WiggidyW/eve-trading-co-go/client/remotedb/appraisal"
	rus "github.com/WiggidyW/eve-trading-co-go/client/remotedb/appraisal/readuserdata"
	"github.com/WiggidyW/eve-trading-co-go/client/remotedb/purchase/cancel"
	"github.com/WiggidyW/eve-trading-co-go/client/shopqueue"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type CancelPurchaseClient struct {
	UserDataClient  rus.SC_ReadUserDataClient
	DelClient       cancel.SMAC_CancelShopPurchaseClient
	ShopQueueClient shopqueue.ShopQueueClient
}

func (cpp CancelPurchaseClient) Fetch(
	ctx context.Context,
	params CancelPurchaseParams,
) (*CancelPurchaseStatus, error) {
	var status CancelPurchaseStatus = Success

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the shop queue
	chnSQRep := util.NewChanResult[[]string](ctx)
	chnSQRepSend, chnSQRepRecv := chnSQRep.Split()
	go cpp.fetchShopQueue(ctx, chnSQRepSend)

	// fetch the user data
	userData, err := cpp.fetchUserData(ctx, params.CharacterId)
	if err != nil {
		return nil, err
	}

	// check if the user can cancel the purchase
	status = checkUserData(*userData, params)
	if status != Success {
		return &status, nil
	}

	// wait for the shop queue
	shopQueue, err := chnSQRepRecv.Recv()
	if err != nil {
		return nil, err
	}

	// check if the appraisal is in the shop queue
	if !appraisalInQueue(shopQueue, params.AppraisalCode) {
		status = PurchaseNotActive
		return &status, nil
	}

	// status = Success
	return &status, nil
}

func (cpp CancelPurchaseClient) fetchUserData(
	ctx context.Context,
	characterId int32,
) (*rdba.UserData, error) {
	if userDataRep, err := cpp.UserDataClient.Fetch(
		ctx,
		rus.ReadUserDataParams{CharacterId: characterId},
	); err != nil {
		return nil, err
	} else {
		userDataVal := userDataRep.Data()
		return &userDataVal, nil
	}
}

func (cpp CancelPurchaseClient) fetchShopQueue(
	ctx context.Context,
	chnSQRepSend util.ChanSendResult[[]string],
) error {
	if sqRep, err := cpp.ShopQueueClient.Fetch(
		ctx,
		shopqueue.ShopQueueParams{},
	); err != nil {
		return chnSQRepSend.SendErr(err)
	} else {
		return chnSQRepSend.SendOk(sqRep.ParsedShopQueue)
	}
}

func appraisalInQueue(shopQueue []string, appraisalCode string) bool {
	for _, code := range shopQueue {
		if code == appraisalCode {
			return true
		}
	}
	return false
}

func cooldownLimited(userData rdba.UserData, cooldown time.Duration) bool {
	return time.Now().Before(userData.CancelledPurchase.Add(cooldown))
}

func characterHasAppraisal(userData rdba.UserData, appraisalCode string) bool {
	for _, code := range userData.ShopAppraisals {
		if code == appraisalCode {
			return true
		}
	}
	return false
}

func checkUserData(
	userData rdba.UserData,
	params CancelPurchaseParams,
) CancelPurchaseStatus {
	cooldownLimited := cooldownLimited(userData, params.Cooldown)
	hasAppraisal := characterHasAppraisal(userData, params.AppraisalCode)
	if cooldownLimited && !hasAppraisal {
		return CooldownLimitedAndPurchaseNotFound
	} else if cooldownLimited {
		return CooldownLimited
	} else if !hasAppraisal {
		return PurchaseNotFound
	} else {
		return Success
	}
}
