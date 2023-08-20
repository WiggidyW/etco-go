package admin

import (
	"context"
	// "github.com/WiggidyW/weve-esi/client/authing"
	"github.com/WiggidyW/weve-esi/client/configure/authhashset/reader"
)

// type AuthingAdminReadClient = authing.AuthingClient[ // "Read" in alrParams
// 	AdminReadParams,
// 	AuthList,
// 	AdminReadClient,
// ]

type AdminReadParams struct {
	refreshToken string
	domain       string
}

func (arcf AdminReadParams) AuthRefreshToken() string {
	return arcf.refreshToken
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
		reader.AuthHashSetReaderParams{ObjectName: params.domain},
	)
	if err != nil {
		return nil, err
	}
	authList := fromHashSet(authHashSet.Data())
	return &authList, nil
}
