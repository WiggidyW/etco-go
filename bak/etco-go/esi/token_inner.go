package esi

import (
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/error/esierror"
)

type EsiAuthRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

var (
	EsiAuthCorp *RefreshTokenAndApp = &RefreshTokenAndApp{
		RefreshToken: build.CORPORATION_WEB_REFRESH_TOKEN,
		App:          EsiAppCorp,
	}
	EsiAuthStructureInfo *RefreshTokenAndApp = &RefreshTokenAndApp{
		RefreshToken: build.STRUCTURE_INFO_WEB_REFRESH_TOKEN,
		App:          EsiAppStructureInfo,
	}
)

func EsiAuthMarkets(refreshToken string) *RefreshTokenAndApp {
	return &RefreshTokenAndApp{
		RefreshToken: refreshToken,
		App:          EsiAppMarkets,
	}
}

type RefreshTokenAndApp struct {
	RefreshToken string
	App          EsiApp
}

func newEsiAccessToken(
	refreshToken string,
	app EsiApp,
	rep EsiAuthRefreshResponse,
) (
	accessToken string,
	expires time.Time,
	err error,
) {
	expires = time.Now().Add(
		time.Duration(rep.ExpiresIn*2/3) * time.Second,
	)
	accessToken = rep.AccessToken
	if refreshToken != rep.RefreshToken {
		err = esierror.AuthRefreshMismatch{App: app.String()}
	}
	return accessToken, expires, err
}
