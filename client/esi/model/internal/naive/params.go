package naive

type NaiveParams[P UrlParams] struct {
	UrlParams  P
	AuthParams *AuthParams
}

func (np NaiveParams[P]) ShouldAuth() bool {
	return np.AuthParams != nil && np.AuthParams.Auth == nil
}

func (np NaiveParams[P]) Auth() *string {
	if np.AuthParams == nil {
		return nil
	} else {
		return np.AuthParams.Auth
	}
}

type AuthParams struct {
	Token string
	Auth  *string
}

type UrlParams interface {
	Url() string
	Method() string
}
