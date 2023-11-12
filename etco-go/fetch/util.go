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

func CalcExpires(
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

func CalcExpiresOptional(
	expires time.Time,
	minExpiresIn *time.Duration,
) time.Time {
	if minExpiresIn == nil {
		return expires
	} else {
		return CalcExpires(expires, *minExpiresIn)
	}
}
