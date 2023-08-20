package readcodes

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
	rdb "github.com/WiggidyW/weve-esi/client/remotedb/internal"
)

func GetCharacterAppraisalCodes(
	ctx context.Context,
	rdbc *rdb.RemoteDBClient,
	characterId int32,
) (*a.CharacterAppraisalCodes, error) {
	fc, err := rdbc.Client(ctx)
	if err != nil {
		return nil, err
	}

	ref := a.CharacterRef(fc, characterId)
	dataTo := &a.CharacterAppraisalCodes{}

	doc, err := ref.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return dataTo, nil
		} else {
			return nil, err
		}
	}

	if err := doc.DataTo(dataTo); err != nil {
		return nil, err
	}

	return dataTo, nil
}
