package service

import (
	"context"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/proto/protoerr"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/staticdb"
)

func (Service) AllShopLocations(
	ctx context.Context,
	_ *proto.EmptyRequest,
) (
	rep *proto.AllShopLocationsResponse,
	err error,
) {
	shopLocationInfos := staticdb.UNSAFE_GetShopLocationInfos()
	numLocations := len(shopLocationInfos)

	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(numLocations)
	rep = &proto.AllShopLocationsResponse{
		Locations: make([]*proto.ShopLocationInfo, numLocations),
	}

	// fetch location infos
	x, cancel := x.WithCancel()
	defer cancel()
	chnInfo :=
		expirable.NewChanResult[*proto.LocationInfo](x.Ctx(), numLocations, 0)
	for locationId := range shopLocationInfos {
		go expirable.P3Transceive(
			chnInfo,
			x, r, locationId,
			esi.ProtoGetLocationInfo,
		)
	}

	// recv location infos
	for i := 0; i < numLocations; i++ {
		locationInfo, _, err := chnInfo.RecvExp()
		if err != nil {
			rep.Error = protoerr.ErrToProto(err)
			rep.Strs = r.Finish()
			return rep, nil
		} else {
			rep.Locations[i] = &proto.ShopLocationInfo{
				LocationInfo: locationInfo,
				TaxRate:      shopLocationInfos[locationInfo.LocationId].TaxRate,
			}
		}
	}

	rep.Strs = r.Finish()
	return rep, nil
}

func (Service) Locations(
	ctx context.Context,
	req *proto.LocationsRequest,
) (
	rep *proto.LocationsResponse,
	err error,
) {
	numLocations := len(req.Locations)

	x := cache.NewContext(ctx)
	r := protoregistry.NewProtoRegistry(numLocations)
	rep = &proto.LocationsResponse{
		Locations: make([]*proto.LocationInfo, numLocations),
	}

	// fetch location infos
	x, cancel := x.WithCancel()
	defer cancel()
	chnInfo :=
		expirable.NewChanResult[*proto.LocationInfo](x.Ctx(), numLocations, 0)
	for _, locationId := range req.Locations {
		go expirable.P3Transceive(
			chnInfo,
			x, r, locationId,
			esi.ProtoGetLocationInfo,
		)
	}

	// recv location infos
	for i := 0; i < numLocations; i++ {
		rep.Locations[i], _, err = chnInfo.RecvExp()
		if err != nil {
			rep.Error = protoerr.ErrToProto(err)
			rep.Strs = r.Finish()
			return rep, nil
		}
	}

	rep.Strs = r.Finish()
	return rep, nil
}
