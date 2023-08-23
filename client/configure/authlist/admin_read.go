package admin

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authingfwding"
	a "github.com/WiggidyW/weve-esi/client/authingfwding/authing"
	"github.com/WiggidyW/weve-esi/client/configure/authhashset/reader"
)

type A_AdminReadClient = a.AuthingClient[
	authingfwding.WithAuthableParams[AdminReadParams],
	AdminReadParams,
	AuthList,
	AdminReadClient,
]

type AdminReadParams struct {
	Domain string
}

type AdminReadClient struct {
	inner reader.AuthHashSetReaderClient
}

func (arc AdminReadClient) Fetch(
	ctx context.Context,
	params AdminReadParams,
) (*AuthList, error) {
	authHashSet, err := arc.inner.Fetch(
		ctx,
		reader.AuthHashSetReaderParams{ObjectName: params.Domain},
	)
	if err != nil {
		return nil, err
	}
	authList := fromHashSet(authHashSet.Data())
	return &authList, nil
}
