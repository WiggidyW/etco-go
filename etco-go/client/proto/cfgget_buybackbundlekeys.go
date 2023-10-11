package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/bucket"
)

type CfgGetBuybackBundleKeysParams struct{}

type CfgGetBuybackBundleKeysClient struct {
	webBundleKeysClient bucket.SC_WebBuybackBundleKeysClient
}

func NewCfgGetBuybackBundleKeysClient(
	webBundleKeysClient bucket.SC_WebBuybackBundleKeysClient,
) CfgGetBuybackBundleKeysClient {
	return CfgGetBuybackBundleKeysClient{webBundleKeysClient}
}

func (gbbkc CfgGetBuybackBundleKeysClient) Fetch(
	ctx context.Context,
	params CfgGetBuybackBundleKeysParams,
) (
	rep *[]string,
	err error,
) {
	webBundleKeysRep, err := gbbkc.webBundleKeysClient.Fetch(
		ctx,
		bucket.WebBuybackBundleKeysParams{},
	)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(webBundleKeysRep.Data()))
	for key := range webBundleKeysRep.Data() {
		keys = append(keys, key)
	}

	return &keys, nil
}
