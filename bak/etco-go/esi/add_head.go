package esi

import (
	"encoding/base64"
	"fmt"
	"net/http"

	build "github.com/WiggidyW/etco-go/buildconstants"
)

const JSON_CONTENT_TYPE = "application/json"
const WWW_CONTENT_TYPE = "application/x-www-form-urlencoded"
const LOGIN_HOST = "login.eveonline.com"

func addHeadUserAgent(req *http.Request) {
	req.Header.Add("X-User-Agent", build.ESI_USER_AGENT)
}

func addHeadJsonContentType(req *http.Request) {
	req.Header.Add("Content-Type", JSON_CONTENT_TYPE)
}

func addHeadWwwContentType(req *http.Request) {
	req.Header.Add("Content-Type", WWW_CONTENT_TYPE)
}

func addHeadBearerAuth(req *http.Request, auth string) {
	if auth != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth))
	}
}

func addHeadBasicAuth(req *http.Request, app EsiApp) {
	clientId, clientSecret := app.ClientIdAndSecret()
	basic_auth := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", clientId, clientSecret),
	))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basic_auth))
}

func addHeadLoginHost(req *http.Request) {
	req.Header.Add("Host", LOGIN_HOST)
}
