package service

import (
	"context"

	built "github.com/WiggidyW/etco-go/builtinconstants"
	"github.com/WiggidyW/etco-go/proto"
)

func (Service) AllAssetFlags(
	_ context.Context,
	_ *proto.EmptyRequest,
) (
	rep *proto.AssetFlagsResponse,
	_ error,
) {
	return &proto.AssetFlagsResponse{Flags: built.ASSET_FLAGS}, nil
}
