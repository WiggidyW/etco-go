package purchase

import (
	"context"
	"time"

	"github.com/WiggidyW/chanresult"

	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	"github.com/WiggidyW/etco-go/client/shopqueue"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

type CancelPurchaseParams struct {
	AppraisalCode string
	CharacterId   int32
	Cooldown      time.Duration // time to wait before allowing character to cancel a purchase
}

type CancelPurchaseClient struct {
	userDataClient  rdbc.SC_ReadUserDataClient
	cancelClient    rdbc.SMAC_CancelPurchaseClient
	shopQueueClient shopqueue.ShopQueueClient
}

func NewCancelPurchaseClient(
	userDataClient rdbc.SC_ReadUserDataClient,
	cancelClient rdbc.SMAC_CancelPurchaseClient,
	shopQueueClient shopqueue.ShopQueueClient,
) CancelPurchaseClient {
	return CancelPurchaseClient{
		userDataClient,
		cancelClient,
		shopQueueClient,
	}
}

func (cpp CancelPurchaseClient) Fetch(
	ctx context.Context,
	params CancelPurchaseParams,
) (*CancelPurchaseStatus, error) {
	var status CancelPurchaseStatus = CPS_Success

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the shop queue
	chnSQRepSend, chnSQRepRecv := chanresult.
		NewChanResult[[]string](ctx, 0, 0).Split()
	go cpp.fetchShopQueue(ctx, chnSQRepSend)

	// fetch the user data
	userData, err := cpp.fetchUserData(ctx, params.CharacterId)
	if err != nil {
		return nil, err
	}

	// check if the user can cancel the purchase
	status = checkUserData(*userData, params)
	if status != CPS_Success {
		return &status, nil
	}

	// wait for the shop queue
	shopQueue, err := chnSQRepRecv.Recv()
	if err != nil {
		return nil, err
	}

	// check if the appraisal is in the shop queue
	if !appraisalInQueue(shopQueue, params.AppraisalCode) {
		status = CPS_PurchaseNotActive
		return &status, nil
	}

	// cancel the purchase
	if _, err := cpp.cancelClient.Fetch(
		ctx,
		rdbc.CancelPurchaseParams{
			CharacterId:   params.CharacterId,
			AppraisalCode: params.AppraisalCode,
		},
	); err != nil {
		return nil, err
	} else {
		// status = Success
		return &status, nil
	}
}

func (cpp CancelPurchaseClient) fetchUserData(
	ctx context.Context,
	characterId int32,
) (*rdb.UserData, error) {
	if userDataRep, err := cpp.userDataClient.Fetch(
		ctx,
		rdbc.ReadUserDataParams{CharacterId: characterId},
	); err != nil {
		return nil, err
	} else {
		userDataVal := userDataRep.Data()
		return &userDataVal, nil
	}
}

func (cpp CancelPurchaseClient) fetchShopQueue(
	ctx context.Context,
	chnSQRepSend chanresult.ChanSendResult[[]string],
) error {
	if sqRep, err := cpp.shopQueueClient.Fetch(
		ctx,
		shopqueue.ShopQueueParams{},
	); err != nil {
		return chnSQRepSend.SendErr(err)
	} else {
		return chnSQRepSend.SendOk(sqRep.ParsedShopQueue)
	}
}

func checkUserData(
	userData rdb.UserData,
	params CancelPurchaseParams,
) CancelPurchaseStatus {
	cooldownLimited := cancelCooldownLimited(userData, params.Cooldown)
	hasAppraisal := characterHasAppraisal(userData, params.AppraisalCode)
	if cooldownLimited && !hasAppraisal {
		return CPS_CooldownLimitedAndPurchaseNotFound
	} else if cooldownLimited {
		return CPS_CooldownLimited
	} else if !hasAppraisal {
		return CPS_PurchaseNotFound
	} else {
		return CPS_Success
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

func cancelCooldownLimited(
	userData rdb.UserData,
	cooldown time.Duration,
) bool {
	return time.Now().Before(userData.CancelledPurchase.Add(cooldown))
}

func characterHasAppraisal(userData rdb.UserData, appraisalCode string) bool {
	for _, code := range userData.ShopAppraisals {
		if code == appraisalCode {
			return true
		}
	}
	return false
}
