package model

import "fmt"

const (
	BASE_URL   = "https://esi.evetech.net/latest"
	DATASOURCE = "tranquility"
)

func AddQueryString(url string, key string, val *string) string {
	if val == nil {
		return url
	} else {
		return fmt.Sprintf("%s&%s=%s", url, key, *val)
	}
}

func AddQueryInt32(url string, key string, val *int32) string {
	if val == nil {
		return url
	} else {
		return fmt.Sprintf("%s&%s=%d", url, key, *val)
	}
}

func AddQueryInt(url string, key string, val *int) string {
	if val == nil {
		return url
	} else {
		return fmt.Sprintf("%s&%s=%d", url, key, *val)
	}
}
