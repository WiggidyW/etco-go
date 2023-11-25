package esi

import (
	"errors"
	"net/http"
	"time"

	"github.com/WiggidyW/etco-go/error/esierror"
)

const (
	ESI_NUM_RETRIES int           = 5
	ESI_RETRY_SLEEP time.Duration = 50 * time.Millisecond
)

func esiShouldRetry(err error) bool {
	var statusErr esierror.StatusError
	return errors.As(err, &statusErr) && esiShouldRetryInner(statusErr)
}

func esiShouldRetryInner(err esierror.StatusError) (retry bool) {
	retry = err.Code == http.StatusBadGateway ||
		err.Code == http.StatusGatewayTimeout
	if retry {
		time.Sleep(ESI_RETRY_SLEEP)
	}
	return retry
}
