package raw_

import (
	"github.com/WiggidyW/etco-go/cache"
)

type EsiAuthResponse struct {
	AccessToken string `json:"access_token"`
}

type EsiAuthResponseWithRefresh = cache.RawAuth
