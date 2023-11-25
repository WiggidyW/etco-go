package esi

import (
	"strings"
	"time"

	"github.com/WiggidyW/etco-go/error/esierror"
)

const (
	// https://github.com/esi/esi-issues/issues/636
	CI_LIMITED_CODE  int           = 520
	CI_LIMITED_MSG   string        = "ConStopSpamming"
	CI_LIMITED_INTVL time.Duration = 10 * time.Second
	CI_LIMITED_REQS  int           = 20
	// max: 20 reqs per 10 secs
)

var (
	ciRateLimiter chan struct{}
)

func init() {
	ciRateLimiter = make(chan struct{}, CI_LIMITED_REQS)
	for i := 0; i < CI_LIMITED_REQS; i++ {
		ciRateLimiter <- struct{}{}
	}
}

func ciRateLimiterStart() {
	<-ciRateLimiter
}

func ciRateLimiterDone() {
	time.Sleep(CI_LIMITED_INTVL)
	ciRateLimiter <- struct{}{}
}

func rateLimited(err esierror.StatusError) (retry bool) {
	retry = err.Code == CI_LIMITED_CODE &&
		strings.Contains(err.EsiText, CI_LIMITED_MSG)
	if retry {
		time.Sleep(CI_LIMITED_INTVL)
	}
	return retry
}
