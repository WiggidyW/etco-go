package bucket

import (
	"context"
)

type AuthListWriterParams struct {
	AuthDomain string
	AuthList   AuthList
}

type AuthListWriterClient struct {
	ahsWriterClient SAC_AuthHashSetWriterClient
}

func NewAuthListWriterClient(
	inner SAC_AuthHashSetWriterClient,
) AuthListWriterClient {
	return AuthListWriterClient{inner}
}

func (ahsrc AuthListWriterClient) Fetch(
	ctx context.Context,
	params AuthListWriterParams,
) (
	rep *struct{},
	err error,
) {
	return ahsrc.ahsWriterClient.Fetch(
		ctx,
		AuthHashSetWriterParams{
			AuthDomain:  params.AuthDomain,
			AuthHashSet: authListToAHS(params.AuthList),
		},
	)
}
