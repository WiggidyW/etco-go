package raw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/error/esierror"
)

const AUTH_URL = "https://login.eveonline.com/v2/oauth/token"
const JWKS_URL = "https://login.eveonline.com/oauth/jwks"

type RawClient struct {
	HttpClient   *http.Client
	UserAgent    string
	ClientId     string
	ClientSecret string
}

func FetchModel[M any](
	rc RawClient,
	ctx context.Context,
	url string,
	method string,
	auth *string,
	model *M,
	// etag *string,
) (*cache.ExpirableData[M], error) {
	// build the request
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, esierror.RequestParamsError{Err: err}
	}
	addHeadUserAgent(req, rc.UserAgent)
	addHeadJsonContentType(req)
	addHeadBearerAuth(req, auth)
	// addHeadEtag(req, etag)

	// fetch the response
	rep, err := rc.HttpClient.Do(req)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return nil, esierror.HttpError{Err: err}
	}

	// if it's not 200, return an error
	if rep.StatusCode != http.StatusOK {
		return nil, esierror.NewStatusError(rep)
	}

	// decode the body
	err = json.NewDecoder(rep.Body).Decode(model)
	if err != nil {
		return nil, esierror.MalformedResponseBody{Err: fmt.Errorf(
			"error decoding response body as json: %w",
			err,
		)}
	}

	// parse the response headers
	expires, err := parseHeadExpires(rep)
	if err != nil {
		return nil, esierror.MalformedResponseHeaders{Err: err}
	}

	output := cache.NewExpirableData[M](*model, expires)
	return &output, nil
}

func (rc RawClient) FetchHead(
	ctx context.Context,
	url string,
	auth *string,
) (*cache.ExpirableData[int], error) {
	// build the request
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return nil, esierror.RequestParamsError{Err: err}
	}
	addHeadUserAgent(req, rc.UserAgent)
	addHeadBearerAuth(req, auth)

	// fetch the response
	rep, err := rc.HttpClient.Do(req)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return nil, esierror.HttpError{Err: err}
	}

	// if it's not 200, return an error
	if rep.StatusCode != http.StatusOK {
		return nil, esierror.NewStatusError(rep)
	}

	// parse the response headers
	expires, err := parseHeadExpires(rep)
	if err != nil {
		return nil, esierror.MalformedResponseHeaders{Err: err}
	}
	pages, err := parseHeadPages(rep)
	if err != nil {
		return nil, esierror.MalformedResponseHeaders{Err: err}
	}

	output := cache.NewExpirableData[int](pages, expires)
	return &output, nil
}

func (rc RawClient) FetchJWKS(
	ctx context.Context,
) (*cache.ExpirableData[[]byte], error) {
	// build the request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		JWKS_URL,
		nil,
	)
	if err != nil {
		return nil, esierror.RequestParamsError{Err: err}
	}

	// fetch the response
	rep, err := rc.HttpClient.Do(req)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return nil, esierror.HttpError{Err: err}
	}

	// if it's not 200, return an error
	if rep.StatusCode != http.StatusOK {
		return nil, esierror.NewStatusError(rep)
	}

	// parse the response headers
	expires, err := parseHeadExpires(rep)
	if err != nil {
		return nil, esierror.MalformedResponseHeaders{Err: err}
	}

	// simply read the body
	// TODO: don't just use 667 (it was the length of the response body when I tested it 8/8/23)
	buf := bytes.NewBuffer(make([]byte, 0, 667))
	_, err = buf.ReadFrom(rep.Body)
	if err != nil {
		return nil, esierror.MalformedResponseBody{Err: fmt.Errorf(
			"error reading response body: %w",
			err,
		)}
	}

	output := cache.NewExpirableData[[]byte](buf.Bytes(), expires)
	return &output, nil
}

func (rc RawClient) FetchAuthWithRefresh(
	ctx context.Context,
	token string,
) (*EsiAuthResponseWithRefresh, error) {
	esiAuthRepWithRefresh := &EsiAuthResponseWithRefresh{}
	err := fetchAuthInner(rc, ctx, token, esiAuthRepWithRefresh)
	if err != nil {
		return nil, err
	}
	return esiAuthRepWithRefresh, nil
}

func (rc RawClient) FetchAuth(
	ctx context.Context,
	token string,
) (*EsiAuthResponse, error) {
	esiAuthRep := &EsiAuthResponse{}
	err := fetchAuthInner(rc, ctx, token, esiAuthRep)
	if err != nil {
		return nil, err
	}
	return esiAuthRep, nil
}

// makes a request to the auth endpoint and encodes the response body into rep
func fetchAuthInner[A any](
	rc RawClient,
	ctx context.Context,
	token string,
	data *A,
) error {
	esiAuthRep := &EsiAuthResponse{}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		AUTH_URL,
		bytes.NewBuffer([]byte(fmt.Sprintf(
			`grant_type=refresh_token&refresh_token=%s`,
			url.QueryEscape(token),
		))),
	)
	if err != nil {
		return esierror.RequestParamsError{Err: err}
	}

	addHeadUserAgent(req, rc.UserAgent)
	addHeadWwwContentType(req)
	addHeadBasicAuth(req, rc.ClientId, rc.ClientSecret)
	addHeadLoginHost(req)

	rep, err := rc.HttpClient.Do(req)
	if err != nil {
		return esierror.HttpError{Err: err}
	}
	defer rep.Body.Close()

	// if it's not 200, return an error
	if rep.StatusCode != http.StatusOK {
		return esierror.NewStatusError(rep)
	}

	err = json.NewDecoder(rep.Body).Decode(esiAuthRep)
	if err != nil {
		return esierror.MalformedResponseBody{Err: fmt.Errorf(
			"error decoding response body as json: %w",
			err,
		)}
	}

	return nil
}