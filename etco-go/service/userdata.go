package service

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/auth"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/remotedb"
)

// returns true if (admin) or (user + same character)
func userAndSameCharacterOrAdmin(
	x cache.Context,
	refreshToken string,
	characterId int32,
) (
	authorized bool,
	err error,
) {
	var rep auth.AuthResponse
	rep, _, err = auth.ProtoUserAuthorized(x, refreshToken)
	if err != nil {
		return false, err
	}
	authorized = rep.AdminStatus == auth.IsAdmin ||
		rep.CharacterId == characterId
	if !authorized {
		rep, _, err = auth.ProtoAdminAuthorized(x, refreshToken)
		if err == nil {
			authorized = rep.Authorized
		}
	}
	return authorized, err
}

func authorizedGetUserDataField[F any](
	x cache.Context,
	req *proto.UserDataRequest,
	getUserDataField func(cache.Context, int32) (F, time.Time, error),
) (
	authorized bool,
	empty F,
	errResponse *proto.ErrorResponse,
) {
	x, cancel := x.WithCancel()
	defer cancel()

	// fetch the field in a goroutine
	chnField := expirable.NewChanResult[F](x.Ctx(), 1, 0)
	go expirable.P2Transceive(
		chnField,
		x, req.CharacterId,
		getUserDataField,
	)

	// check if user is authorized
	var err error
	authorized, err =
		userAndSameCharacterOrAdmin(x, req.RefreshToken, req.CharacterId)
	if !authorized || err != nil {
		return false, empty, protoerr.ErrToProto(err)
	}

	// recv the field
	var field F
	field, _, err = chnField.RecvExp()
	if err != nil {
		errResponse = protoerr.ErrToProto(err)
	}
	return authorized, field, errResponse
}

func authorizedGetUserDataTimestamp(
	x cache.Context,
	req *proto.UserDataRequest,
	getUserDataTimestamp func(
		cache.Context,
		int32,
	) (*time.Time, time.Time, error),
) (
	authorized bool,
	timestamp int64,
	errResponse *proto.ErrorResponse,
) {
	var field *time.Time
	authorized, field, errResponse =
		authorizedGetUserDataField(x, req, getUserDataTimestamp)
	if errResponse == nil && field != nil {
		timestamp = field.Unix()
	}
	return authorized, timestamp, errResponse
}

func (Service) UserBuybackAppraisalCodes(
	ctx context.Context,
	req *proto.UserDataRequest,
) (
	rep *proto.UserAppraisalCodesResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.UserAppraisalCodesResponse{}
	rep.Authorized, rep.Codes, rep.Error = authorizedGetUserDataField(
		x,
		req,
		remotedb.GetUserBuybackAppraisalCodes,
	)
	return rep, nil
}

func (Service) UserShopAppraisalCodes(
	ctx context.Context,
	req *proto.UserDataRequest,
) (
	rep *proto.UserAppraisalCodesResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.UserAppraisalCodesResponse{}
	rep.Authorized, rep.Codes, rep.Error = authorizedGetUserDataField(
		x,
		req,
		remotedb.GetUserShopAppraisalCodes,
	)
	return rep, nil
}

func (Service) UserHaulAppraisalCodes(
	ctx context.Context,
	req *proto.UserDataRequest,
) (
	rep *proto.UserAppraisalCodesResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.UserAppraisalCodesResponse{}
	rep.Authorized, rep.Codes, rep.Error = authorizedGetUserDataField(
		x,
		req,
		remotedb.GetUserHaulAppraisalCodes,
	)
	return rep, nil
}

func (Service) UserCancelledPurchase(
	ctx context.Context,
	req *proto.UserDataRequest,
) (
	rep *proto.UserTimePurchaseResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.UserTimePurchaseResponse{}
	rep.Authorized, rep.Time, rep.Error = authorizedGetUserDataTimestamp(
		x,
		req,
		remotedb.GetUserCancelledPurchase,
	)
	return rep, nil
}

func (Service) UserMadePurchase(
	ctx context.Context,
	req *proto.UserDataRequest,
) (
	rep *proto.UserTimePurchaseResponse,
	err error,
) {
	x := cache.NewContext(ctx)
	rep = &proto.UserTimePurchaseResponse{}
	rep.Authorized, rep.Time, rep.Error = authorizedGetUserDataTimestamp(
		x,
		req,
		remotedb.GetUserMadePurchase,
	)
	return rep, nil
}
