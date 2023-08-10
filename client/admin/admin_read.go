package admin

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authing"
)

type AuthingAdminReadClient = authing.AuthingClient[ // "Read" in alrParams
	AdminReadParams,
	AuthList,
	AdminReadClient,
]

type AdminReadParams struct {
	refreshToken string
	domain       string
}

func (arcf AdminReadParams) AuthRefreshToken() string {
	return arcf.refreshToken
}

type AdminReadClient struct {
	inner authing.AuthHashSetReaderClient
}

func (arc AdminReadClient) Fetch(
	ctx context.Context,
	params AdminReadParams,
) (*AuthList, error) {
	authHashSet, err := arc.inner.Fetch(
		ctx,
		authing.AuthHashSetReaderParams(params.domain),
	)
	if err != nil {
		return nil, err
	}
	authList := fromHashSet(authHashSet.Data())
	return &authList, nil
}
