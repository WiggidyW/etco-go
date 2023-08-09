package modelclient

import "fmt"

const (
	BASE_URL   = "https://esi.evetech.net/latest"
	DATASOURCE = "tranquility"
)

func addQueryString(url string, key string, val *string) string {
	if val == nil {
		return url
	} else {
		return fmt.Sprintf("%s&%s=%s", url, key, *val)
	}
}

func addQueryInt32(url string, key string, val *int32) string {
	if val == nil {
		return url
	} else {
		return fmt.Sprintf("%s&%s=%d", url, key, *val)
	}
}
