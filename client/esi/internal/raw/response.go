package raw

type EsiAuthResponse struct {
	AccessToken string `json:"access_token"`
}

type EsiAuthResponseWithRefresh struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
