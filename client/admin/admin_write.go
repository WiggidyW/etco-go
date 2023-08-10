package admin

import (
	"context"

	"github.com/WiggidyW/weve-esi/client/authing"
)

type AuthingAdminWriteClient = authing.AuthingClient[ // "Write" in alrParams
	AdminWriteParams,
	struct{},
	AdminWriteClient,
]

type AdminWriteParams struct {
	refreshToken string
	domain       string
	authList     AuthList
}

func (awcf AdminWriteParams) AuthRefreshToken() string {
	return awcf.refreshToken
}

type AdminWriteClient struct {
	inner authing.AuthHashSetWriterClient
}

func (awc AdminWriteClient) Fetch(
	ctx context.Context,
	params AdminWriteParams,
) (*struct{}, error) {
	return awc.inner.Fetch(
		ctx,
		authing.AuthHashSetWriterParams{
			Key: params.domain,
			Val: params.authList.toHashSet(),
		},
	)
}
