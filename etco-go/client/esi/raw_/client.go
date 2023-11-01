package raw_

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/error/esierror"
	"github.com/WiggidyW/etco-go/logger"
)

const (
	AUTH_URL string = "https://login.eveonline.com/v2/oauth/token"
	JWKS_URL string = "https://login.eveonline.com/oauth/jwks"

	NUM_RETRIES int           = 5
	RETRY_WAIT  time.Duration = 500 * time.Millisecond
)

type RawClient struct {
	HttpClient   *http.Client
	UserAgent    string
	ClientId     string
	ClientSecret string
	Namespace    string
	TokenCache   *cache.TokenCache
}

func NewUnauthenticatedRawClient(httpClient *http.Client) RawClient {
	return RawClient{
		HttpClient: httpClient,
		UserAgent:  build.ESI_USER_AGENT,
		// ClientId: "",
		// ClientSecret: "",
		// Namespace: "",
		// TokenCache: nil,
	}
}

func NewCorpRawClient(
	httpClient *http.Client,
	cCache cache.SharedClientCache,
) RawClient {
	return RawClient{
		HttpClient:   httpClient,
		UserAgent:    build.ESI_USER_AGENT,
		ClientId:     build.ESI_CORP_CLIENT_ID,
		ClientSecret: build.ESI_CORP_CLIENT_SECRET,
		Namespace:    "corp",
		TokenCache:   cache.NewTokenCache(cCache),
	}
}

func NewMarketsRawClient(
	httpClient *http.Client,
	cCache cache.SharedClientCache,
) RawClient {
	return RawClient{
		HttpClient:   httpClient,
		UserAgent:    build.ESI_USER_AGENT,
		ClientId:     build.ESI_MARKETS_CLIENT_ID,
		ClientSecret: build.ESI_MARKETS_CLIENT_SECRET,
		Namespace:    "markets",
		TokenCache:   cache.NewTokenCache(cCache),
	}
}

func NewStructureInfoRawClient(
	httpClient *http.Client,
	cCache cache.SharedClientCache,
) RawClient {
	return RawClient{
		HttpClient:   httpClient,
		UserAgent:    build.ESI_USER_AGENT,
		ClientId:     build.ESI_STRUCTURE_INFO_CLIENT_ID,
		ClientSecret: build.ESI_STRUCTURE_INFO_CLIENT_SECRET,
		Namespace:    "struc",
		TokenCache:   cache.NewTokenCache(cCache),
	}
}

func NewAuthRawClient(
	httpClient *http.Client,
	cCache cache.SharedClientCache,
) RawClient {
	return RawClient{
		HttpClient:   httpClient,
		UserAgent:    build.ESI_USER_AGENT,
		ClientId:     build.ESI_AUTH_CLIENT_ID,
		ClientSecret: build.ESI_AUTH_CLIENT_SECRET,
		Namespace:    "auth",
		TokenCache:   cache.NewTokenCache(cCache),
	}
}

func FetchModel[M any](
	rc RawClient,
	ctx context.Context,
	url string,
	method string,
	auth *string,
	model *M,
	// etag *string,
) (*cache.ExpirableData[M], error) { return fetchWithRetries1[*cache.ExpirableData[M]](func() (*cache.ExpirableData[M], error) {
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
}, 1)}

func (rc RawClient) FetchHead(
	ctx context.Context,
	url string,
	auth *string,
) (*cache.ExpirableData[int], error) { return fetchWithRetries1[*cache.ExpirableData[int]](func() (*cache.ExpirableData[int], error) {
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
}, 1)}

func (rc RawClient) FetchJWKS(
	ctx context.Context,
) (*cache.ExpirableData[[]byte], error) { return fetchWithRetries1[*cache.ExpirableData[[]byte]](func() (*cache.ExpirableData[[]byte], error) {
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
	addHeadUserAgent(req, rc.UserAgent)

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
}, 1)}

func (rc RawClient) FetchAuthWithRefreshFromCode(
	ctx context.Context,
	code string,
) (*EsiAuthResponseWithRefresh, error) {
	esiAuthRepWithRefresh := &EsiAuthResponseWithRefresh{}
	err := fetchAuthInner(
		rc,
		ctx,
		fetchAuthCodeBody(code),
		esiAuthRepWithRefresh,
	)
	if err != nil {
		return nil, err
	}
	return esiAuthRepWithRefresh, nil
}

func (rc RawClient) FetchAuthWithRefresh(
	ctx context.Context,
	token string,
) (*EsiAuthResponseWithRefresh, error) {
	return fetchCacheAuthInner(rc, ctx, token)
}

func (rc RawClient) FetchAuth(
	ctx context.Context,
	token string,
) (*EsiAuthResponse, error) {
	rep, err := fetchCacheAuthInner(rc, ctx, token)
	if err != nil {
		return nil, err
	}
	return &EsiAuthResponse{AccessToken: rep.AccessToken}, nil
}

func fetchAuthCodeBody(
	code string,
) string {
	return fmt.Sprintf(
		`grant_type=authorization_code&code=%s`,
		url.QueryEscape(code),
	)
}

func fetchAuthTokenBody(
	refreshToken string,
) string {
	return fmt.Sprintf(
		`grant_type=refresh_token&refresh_token=%s`,
		url.QueryEscape(refreshToken),
	)
}

func fetchCacheAuthInner(
	rc RawClient,
	ctx context.Context,
	token string,
) (
	rep *EsiAuthResponseWithRefresh,
	err error,
) {
	cacheKey := fmt.Sprintf("%s-%s", rc.Namespace, token)

	rep, lock, err := rc.TokenCache.GetOrLock(cacheKey)
	if err != nil {
		return nil, err
	} else if rep != nil {
		return rep, nil
	}

	rep = &EsiAuthResponseWithRefresh{}
	expires := time.Now().Add(time.Minute * 10)
	err = fetchAuthInner(rc, ctx, fetchAuthTokenBody(token), rep)
	if err != nil {
		go func() { rc.TokenCache.Unlock(lock) }()
		return nil, err
	}

	go func() {
		logger.Err(rc.TokenCache.Set(cacheKey, *rep, expires, lock))
	}()
	return rep, nil
}

// makes a request to the auth endpoint and encodes the response body into rep
func fetchAuthInner[A any](
	rc RawClient,
	ctx context.Context,
	body string, data *A, //prettier-ignore
) error { return fetchWithRetries0(func() error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		AUTH_URL,
		bytes.NewBuffer([]byte(body)),
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

	err = json.NewDecoder(rep.Body).Decode(data)
	if err != nil {
		return esierror.MalformedResponseBody{Err: fmt.Errorf(
			"error decoding response body as json: %w",
			err,
		)}
	}

	return nil
}, 1)}

func fetchWithRetries0(
	fn func() error,
	attempt int,
) error {
	err := fn()
	if err != nil && attempt < NUM_RETRIES && shouldRetry(err) {
		time.Sleep(RETRY_WAIT)
		return fetchWithRetries0(fn, attempt+1)
	} else {
		return err
	}
}

func fetchWithRetries1[T any](
	fn func() (T, error),
	attempt int,
) (T, error) {
	rep, err := fn()
	if err != nil && attempt < NUM_RETRIES && shouldRetry(err) {
		time.Sleep(RETRY_WAIT)
		return fetchWithRetries1(fn, attempt+1)
	} else {
		return rep, err
	}
}
