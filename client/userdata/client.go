package userdata

import (
	"context"

	"github.com/WiggidyW/chanresult"
	"github.com/WiggidyW/etco-go/client/contracts"
	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	"github.com/WiggidyW/etco-go/client/shopqueue"
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/util"
)

type UserDataClient struct {
	readClient      rdbc.SC_ReadUserDataClient
	shopQueueClient shopqueue.ShopQueueClient
}

func NewUserDataClient(
	readClient rdbc.SC_ReadUserDataClient,
	shopQueueClient shopqueue.ShopQueueClient,
) UserDataClient {
	return UserDataClient{readClient, shopQueueClient}
}

func (udc UserDataClient) Fetch(
	ctx context.Context,
	params UserDataParams,
) (
	userData UserData,
	err error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fetch the raw user data in a separate goroutine
	chnSendRawUserData, chnRecvRawUserData := chanresult.
		NewChanResult[rdb.UserData](ctx, 1, 0).Split()
	go udc.transceiveFetchRawUserData(ctx, params, chnSendRawUserData)

	// fetch the shop queue
	shopQueueRep, err := udc.fetchShopQueue(ctx)
	if err != nil {
		return userData, err
	}

	// get a purchase queue hash set from the shop queue
	sqHashSet := shopQueueRep.ShopQueueHashSet()

	// wait for the raw user data
	rawUserData, err := chnRecvRawUserData.Recv()
	if err != nil {
		return userData, err
	} else {
		userData.MadePurchase = rawUserData.MadePurchase
		userData.CancelledPurchase = rawUserData.CancelledPurchase
	}

	// populate the buyback appraisals
	userData.BuybackAppraisals = make(
		[]BuybackAppraisalStatus,
		0,
		len(rawUserData.BuybackAppraisals),
	)
	for _, code := range rawUserData.BuybackAppraisals {
		addBuybackAppraisal(
			&userData.BuybackAppraisals,
			code,
			shopQueueRep.BuybackContracts,
		)
	}

	// populate the shop appraisals
	userData.ShopAppraisals = make(
		[]ShopAppraisalStatus,
		0,
		len(rawUserData.ShopAppraisals),
	)
	for _, code := range rawUserData.ShopAppraisals {
		addShopAppraisal(
			&userData.ShopAppraisals,
			code,
			shopQueueRep.ShopContracts,
			sqHashSet,
		)
	}

	return userData, nil
}

func (udc UserDataClient) transceiveFetchRawUserData(
	ctx context.Context,
	params UserDataParams,
	chnSend chanresult.ChanSendResult[rdb.UserData],
) error {
	rawUserData, err := udc.fetchRawUserData(ctx, params)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(rawUserData)
	}
}

func (udc UserDataClient) fetchRawUserData(
	ctx context.Context,
	params UserDataParams,
) (
	rawUserData rdb.UserData,
	err error,
) {
	rawUserDataRep, err := udc.readClient.Fetch(
		ctx,
		rdbc.ReadUserDataParams(params),
	)
	if err != nil {
		return rawUserData, err
	} else {
		return rawUserDataRep.Data(), nil
	}
}

// func (udc UserDataClient) transceiveFetchShopQueue(
// 	ctx context.Context,
// 	chnSend chanresult.ChanSendResult[shopqueue.ShopQueueResponse],
// ) error {
// 	shopQueue, err := udc.fetchShopQueue(ctx)
// 	if err != nil {
// 		return chnSend.SendErr(err)
// 	} else {
// 		return chnSend.SendOk(shopQueue)
// 	}
// }

func (udc UserDataClient) fetchShopQueue(ctx context.Context) (
	shopQueue shopqueue.ShopQueueResponse,
	err error,
) {
	shopQueueRep, err := udc.shopQueueClient.Fetch(
		ctx,
		shopqueue.ShopQueueParams{},
	)
	if err != nil {
		return shopQueue, err
	} else {
		return *shopQueueRep, nil
	}
}

func addBuybackAppraisal(
	buybackAppraisals *[]BuybackAppraisalStatus,
	appraisalCode string,
	buybackContracts map[string]contracts.Contract,
) {
	if contract, ok := buybackContracts[appraisalCode]; ok {
		*buybackAppraisals = append(
			*buybackAppraisals,
			BuybackAppraisalStatus{
				Code:     appraisalCode,
				Contract: &contract,
			},
		)
	} else {
		*buybackAppraisals = append(
			*buybackAppraisals,
			BuybackAppraisalStatus{
				Code: appraisalCode,
				// Contract: nil,
			},
		)
	}
}

func addShopAppraisal(
	shopAppraisals *[]ShopAppraisalStatus,
	appraisalCode string,
	shopContracts map[string]contracts.Contract,
	shopSQHashSet util.MapHashSet[string, struct{}],
) {
	if shopSQHashSet.Has(appraisalCode) {
		*shopAppraisals = append(
			*shopAppraisals,
			ShopAppraisalStatus{
				Code:            appraisalCode,
				InPurchaseQueue: true,
				// Contract: nil,
			},
		)
	} else if contract, ok := shopContracts[appraisalCode]; ok {
		*shopAppraisals = append(
			*shopAppraisals,
			ShopAppraisalStatus{
				Code: appraisalCode,
				// InPurchaseQueue: false,
				Contract: &contract,
			},
		)
	} else {
		*shopAppraisals = append(
			*shopAppraisals,
			ShopAppraisalStatus{
				Code: appraisalCode,
				// InPurchaseQueue: false,
				// Contract: nil,
			},
		)
	}
}
