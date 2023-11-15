package proto

import (
	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
)

type CfgGetShopBundleKeysParams struct{}

type CfgGetShopBundleKeysClient struct{}

func NewCfgGetShopBundleKeysClient() CfgGetShopBundleKeysClient {
	return CfgGetShopBundleKeysClient{}
}

func (gsbkc CfgGetShopBundleKeysClient) Fetch(
	x cache.Context,
	params CfgGetShopBundleKeysParams,
) (
	rep *[]string,
	err error,
) {
	var webBundleKeysRep map[string]struct{}
	webBundleKeysRep, _, err = bucket.GetWebShopBundleKeys(x)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(webBundleKeysRep))
	for key := range webBundleKeysRep {
		keys = append(keys, key)
	}

	return &keys, nil
}
