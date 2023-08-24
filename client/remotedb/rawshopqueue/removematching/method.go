package removematching

import (
	"context"

	"cloud.google.com/go/firestore"
	rdb "github.com/WiggidyW/eve-trading-co-go/client/remotedb/internal"
	sq "github.com/WiggidyW/eve-trading-co-go/client/remotedb/rawshopqueue"
)

func SetShopQueueRemoveMatching(
	rdbc *rdb.RemoteDBClient,
	ctx context.Context,
	remove []string,
) error {
	fc, err := rdbc.Client(ctx)
	if err != nil {
		return err
	}

	iRemove := make([]interface{}, len(remove))
	for i, v := range remove {
		iRemove[i] = v
	}

	_, err = fc.Collection(sq.COLLECTION_ID).Doc(sq.DOC_ID).Set(
		ctx,
		map[string]interface{}{
			sq.FIELD_ID: firestore.ArrayRemove(iRemove...),
		},
		firestore.MergeAll,
	)

	return err
}
