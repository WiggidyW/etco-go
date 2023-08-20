package appraisal

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CharacterRef(
	fc *firestore.Client,
	characterId int32,
) *firestore.DocumentRef {
	return fc.Collection(CHARACTERS_COLLECTION_ID).Doc(
		fmt.Sprintf("%d", characterId),
	)
}

// returns exists false if an error occurs before existence can be known
func GetAppraisal(
	rdbc *rdb.RemoteDBClient,
	ctx context.Context,
	appraisalKey string,
	collectionId string,
	dataTo interface{}, // must be a pointer
) (exists bool, err error) {
	fc, err := rdbc.Client(ctx)
	if err != nil {
		return false, err
	}

	doc, err := fc.Collection(collectionId).Doc(appraisalKey).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	if err := doc.DataTo(dataTo); err != nil {
		return true, err
	}

	return true, nil
}
