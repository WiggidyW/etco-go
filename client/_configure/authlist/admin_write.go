package admin

import (
	"context"

	"github.com/WiggidyW/etco-go/client/authingfwding"
	a "github.com/WiggidyW/etco-go/client/authingfwding/authing"
	"github.com/WiggidyW/etco-go/client/configure/authhashset/writer"
)

// type AuthingAdminWriteClient = authing.AuthingClient[ // "Write" in alrParams
// 	AdminWriteParams,
// 	struct{},
// 	AdminWriteClient,
// ]

type A_AdminWriteClient = a.AuthingClient[
	authingfwding.WithAuthableParams[AdminWriteParams],
	AdminWriteParams,
	struct{},
	AdminWriteClient,
]

type AdminWriteClient struct {
	inner writer.AuthHashSetWriterClient
}

func (awc AdminWriteClient) Fetch(
	ctx context.Context,
	params AdminWriteParams,
) (*struct{}, error) {
	return awc.inner.Fetch(
		ctx,
		writer.AuthHashSetWriterParams{
			ObjectName: params.Domain,
			Val:        params.AuthList.toHashSet(),
		},
	)
}

type AdminWriteParams struct {
	Domain   string
	AuthList AuthList
}
