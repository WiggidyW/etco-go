package esi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/error/esierror"
)

const (
	AUTH_URL   string = "https://login.eveonline.com/v2/oauth/token"
	JWKS_URL   string = "https://login.eveonline.com/oauth/jwks"
	BASE_URL   string = "https://esi.evetech.net/latest"
	DATASOURCE string = "tranquility"
)

var (
	client = &http.Client{}
)

func authLogin(
	ctx context.Context,
	accessCode string,
	app EsiApp,
) (
	refreshToken string,
	err error,
) {
	// build the request
	var req *http.Request
	req, err = newRequest(
		ctx,
		http.MethodPost,
		AUTH_URL,
		authLoginBody(accessCode),
	)
	if err != nil {
		return refreshToken, err
	}
	addHeadWwwContentType(req)
	addHeadBasicAuth(req, app)
	addHeadLoginHost(req)

	// fetch the response
	var httpRep *http.Response
	var close func() error
	httpRep, close, err = doRequest(req)
	defer close()
	if err != nil {
		return refreshToken, err
	}

	// decode the body
	var authRep EsiAuthRefreshResponse
	_, err = decode(httpRep.Body, &authRep)

	return authRep.RefreshToken, err
}
func authLoginBody(accessCode string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(fmt.Sprintf(
		`grant_type=authorization_code&code=%s`,
		url.QueryEscape(accessCode),
	)))
}

func authRefresh(
	ctx context.Context,
	refreshToken string,
	app EsiApp,
) (
	accessToken string,
	expires time.Time,
	err error,
) {
	// build the request
	var req *http.Request
	req, err = newRequest(
		ctx,
		http.MethodPost,
		AUTH_URL,
		authRefreshBody(refreshToken),
	)
	if err != nil {
		return accessToken, expires, err
	}
	addHeadWwwContentType(req)
	addHeadBasicAuth(req, app)
	addHeadLoginHost(req)

	// fetch the response
	var httpRep *http.Response
	var close func() error
	httpRep, close, err = doRequest(req)
	defer close()
	if err != nil {
		return accessToken, expires, err
	}

	// decode the body
	var authRep EsiAuthRefreshResponse
	_, err = decode(httpRep.Body, &authRep)
	if err != nil {
		return accessToken, expires, err
	}

	// validate the response
	return newEsiAccessToken(refreshToken, app, authRep)
}
func authRefreshBody(refreshToken string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(fmt.Sprintf(
		`grant_type=refresh_token&refresh_token=%s`,
		url.QueryEscape(refreshToken),
	)))
}

func getModel[M any](
	x cache.Context,
	url string,
	method string,
	auth *RefreshTokenAndApp,
	newModel func() *M,
) (
	model M,
	expires time.Time,
	err error,
) {
	// get accessToken if auth is not nil
	var accessToken string
	if auth != nil {
		accessToken, _, err = accessTokenGet(x, auth.RefreshToken, auth.App)
		if err != nil {
			return model, expires, err
		}
	}

	// build the request
	var req *http.Request
	req, err = newRequest(x.Ctx(), method, url, nil)
	if err != nil {
		return model, expires, err
	}
	addHeadJsonContentType(req)
	addHeadBearerAuth(req, accessToken)

	// fetch the response
	var httpRep *http.Response
	var close func() error
	httpRep, close, err = doRequest(req)
	defer close()
	if err != nil {
		return model, expires, err
	}

	// decode the body
	model, err = decode(httpRep.Body, newRepOrDefault(newModel))
	if err != nil {
		return model, expires, err
	}

	// parse the response headers
	expires, err = parseHeadExpires(httpRep)
	if err != nil {
		return model, expires, err
	}

	return model, expires, nil
}

func getHead(
	x cache.Context,
	url string,
	auth *RefreshTokenAndApp,
) (
	pages int,
	expires time.Time,
	err error,
) {
	// get accessToken if auth is not nil
	var accessToken string
	if auth != nil {
		accessToken, _, err = accessTokenGet(x, auth.RefreshToken, auth.App)
		if err != nil {
			return pages, expires, err
		}
	}

	// build the request
	var req *http.Request
	req, err = newRequest(x.Ctx(), http.MethodHead, url, nil)
	if err != nil {
		return pages, expires, err
	}
	addHeadBearerAuth(req, accessToken)

	// fetch the response
	var httpRep *http.Response
	var close func() error
	httpRep, close, err = doRequest(req)
	defer close()
	if err != nil {
		return pages, expires, err
	}

	// parse the response headers
	expires, err = parseHeadExpires(httpRep)
	if err != nil {
		return pages, expires, err
	}
	pages, err = parseHeadPages(httpRep)
	if err != nil {
		return pages, expires, err
	}

	return pages, expires, nil
}

func getJWKS(
	ctx context.Context,
	buf []byte,
) (
	jwks []byte,
	expires time.Time,
	err error,
) {
	// build the request
	var req *http.Request
	req, err = newRequest(ctx, http.MethodGet, JWKS_URL, nil)
	if err != nil {
		return nil, expires, err
	}

	// fetch the response
	var httpRep *http.Response
	var close func() error
	httpRep, close, err = doRequest(req)
	defer close()
	if err != nil {
		return nil, expires, err
	}

	// simply read the body
	writer := bytes.NewBuffer(buf)
	_, err = writer.ReadFrom(httpRep.Body)
	if err != nil {
		return nil, expires, esierror.MalformedResponseBody{Err: fmt.Errorf(
			"error reading response body: %w",
			err,
		)}
	}

	// parse the response headers
	expires, err = parseHeadExpires(httpRep)
	if err != nil {
		return nil, expires, err
	}

	return writer.Bytes(), expires, nil
}

func newRequest(
	ctx context.Context,
	method, url string,
	body io.Reader,
) (
	req *http.Request,
	err error,
) {
	req, err = http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		err = esierror.RequestParamsError{Err: err}
	} else {
		addHeadUserAgent(req)
	}
	return req, err
}

func doRequest(
	req *http.Request,
) (
	rep *http.Response,
	close func() error,
	err error,
) {
	rep, err = client.Do(req)
	if err != nil {
		close = voidClose
		err = esierror.HttpError{Err: err}
	} else {
		close = rep.Body.Close
		if rep.StatusCode != http.StatusOK {
			err = esierror.NewStatusError(rep)
		}
	}
	return rep, close, err
}
func voidClose() error { return nil }

func newRepOrDefault[REP any](
	newRep func() *REP,
) *REP {
	if newRep != nil {
		return newRep()
	} else {
		return new(REP)
	}
}

func decode[REP any](
	body io.ReadCloser,
	repPtr *REP,
) (
	rep REP,
	err error,
) {
	err = json.NewDecoder(body).Decode(repPtr)
	if err != nil {
		err = esierror.MalformedResponseBody{Err: fmt.Errorf(
			"error decoding response body as json: %w",
			err,
		)}
	} else if repPtr != nil {
		rep = *repPtr
	}
	return rep, err
}
