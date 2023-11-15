package proto

import (
	"github.com/WiggidyW/etco-go/bucket"
	"github.com/WiggidyW/etco-go/cache"
)

type CfgGetBuybackBundleKeysParams struct{}

type CfgGetBuybackBundleKeysClient struct{}

func NewCfgGetBuybackBundleKeysClient() CfgGetBuybackBundleKeysClient {
	return CfgGetBuybackBundleKeysClient{}
}

func (gbbkc CfgGetBuybackBundleKeysClient) Fetch(
	x cache.Context,
	params CfgGetBuybackBundleKeysParams,
) (
	rep *[]string,
	err error,
) {
	var webBundleKeysRep map[string]struct{}
	webBundleKeysRep, _, err = bucket.GetWebBuybackBundleKeys(x)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(webBundleKeysRep))
	for key := range webBundleKeysRep {
		keys = append(keys, key)
	}

	return &keys, nil
}
