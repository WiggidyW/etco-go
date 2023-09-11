package proto

import (
	"context"

	"github.com/WiggidyW/etco-go/client/bucket"
)

type CfgGetShopBundleKeysParams struct{}

type CfgGetShopBundleKeysClient struct {
	webBundleKeysClient bucket.SC_WebShopBundleKeysClient
}

func NewCfgGetShopBundleKeysClient(
	webBundleKeysClient bucket.SC_WebShopBundleKeysClient,
) CfgGetShopBundleKeysClient {
	return CfgGetShopBundleKeysClient{webBundleKeysClient}
}

func (gsbkc CfgGetShopBundleKeysClient) Fetch(
	ctx context.Context,
	params CfgGetShopBundleKeysParams,
) (
	rep *[]string,
	err error,
) {
	webBundleKeysRep, err := gsbkc.webBundleKeysClient.Fetch(
		ctx,
		bucket.WebShopBundleKeysParams{},
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
