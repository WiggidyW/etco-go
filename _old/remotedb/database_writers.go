package db

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/authing"
)

type MakePurchaseParams[S IShopAppraisal[I], I IShopAppraisalItem] struct {
	AppraisalKey string
	Appraisal    S
	CharacterId  int32
}

type MakePurchaseClient[S IShopAppraisal[I], I IShopAppraisalItem] struct {
	*DatabaseClient
}

func (mpc MakePurchaseClient[S, I]) Fetch(
	ctx context.Context,
	params MakePurchaseParams[S, I],
) (*struct{}, error) {
	if err := MakePurchase[S, I](
		mpc.DatabaseClient,
		ctx,
		params.CharacterId,
		params.AppraisalKey,
		params.Appraisal,
	); err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}

type AF_SAC_CancelPurchaseClient = authing.AuthFwdingClient[
	FwdableCancelPurchaseParams,
	CancelPurchaseParams,
	struct{},
	client.StrongMultiAntiCachingClient[
		CancelPurchaseParams,
		struct{},
		CancelPurchaseClient,
	],
]

type FwdableCancelPurchaseParams struct {
	AppraisalKey string
	RefreshToken string
}

func (fcpp FwdableCancelPurchaseParams) AuthRefreshToken() string {
	return fcpp.RefreshToken
}

func (fcpp FwdableCancelPurchaseParams) WithCharacterId(
	CharacterId int32,
) CancelPurchaseParams {
	return CancelPurchaseParams{fcpp.AppraisalKey, CharacterId}
}

type CancelPurchaseParams struct {
	AppraisalKey string
	CharacterId  int32
}

func (cpp CancelPurchaseParams) AntiCacheKeys() []string {
	return []string{
		"shopqueue",
		fmt.Sprintf("charcodes-%d", cpp.CharacterId),
	}
}

type CancelPurchaseClient struct {
	*DatabaseClient
}

func (cpc CancelPurchaseClient) Fetch(
	ctx context.Context,
	params CancelPurchaseParams,
) (*struct{}, error) {
	if err := cpc.CancelPurchase(
		ctx,
		params.CharacterId,
		params.AppraisalKey,
	); err != nil {
		return nil, err
	} else {
		return &struct{}{}, nil
	}
}
