package fetch

import (
	"math"
	"time"
)

var (
	MAX_EXPIRES time.Time = time.Unix(
		math.MaxInt64,
		math.MaxInt64,
	)
)

func CalcExpires(expires1 time.Time, expires2 time.Time) time.Time {
	if expires1.Before(expires2) {
		return expires1
	} else {
		return expires2
	}
}

func CalcExpiresIn(
	expires time.Time,
	minExpiresIn time.Duration,
) time.Time {
	minExpires := time.Now().Add(minExpiresIn)
	if expires.Before(minExpires) {
		return minExpires
	} else {
		return expires
	}
}

func CalcExpiresInOptional(
	expires time.Time,
	minExpiresIn *time.Duration,
) time.Time {
	if minExpiresIn == nil {
		return expires
	} else {
		return CalcExpiresIn(expires, *minExpiresIn)
	}
}
