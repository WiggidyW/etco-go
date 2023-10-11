package bucket

import (
	"context"
)

type AuthListReaderParams struct {
	AuthDomain string
}

type AuthListReaderClient struct {
	ahsReaderClient SC_AuthHashSetReaderClient
}

func NewAuthListReaderClient(
	inner SC_AuthHashSetReaderClient,
) AuthListReaderClient {
	return AuthListReaderClient{inner}
}

func (alrc AuthListReaderClient) Fetch(
	ctx context.Context,
	params AuthListReaderParams,
) (
	rep AuthList,
	err error,
) {
	ahsRep, err := alrc.ahsReaderClient.Fetch(
		ctx,
		AuthHashSetReaderParams(params),
	)
	if err != nil {
		return rep, err
	} else {
		return authListfromAHS(ahsRep.Data()), nil
	}
}
