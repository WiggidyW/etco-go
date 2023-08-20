package read

import (
	"context"

	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
	sq "github.com/WiggidyW/weve-esi/client/remotedb/rawshopqueue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetShopQueue(
	rdbc *rdb.RemoteDBClient,
	ctx context.Context,
) ([]string, error) {
	fc, err := rdbc.Client(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := fc.Collection(sq.COLLECTION_ID).Doc(sq.DOC_ID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return []string{}, nil
		} else {
			return nil, err
		}
	}

	sq := new(sq.ShopQueue)
	if err := doc.DataTo(&sq); err != nil {
		return nil, err
	}

	return sq.ShopQueue, nil
}
