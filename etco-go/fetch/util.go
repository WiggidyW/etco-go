package fetch

import (
	"time"
)

func ExpiresIn(expiresIn time.Duration) *time.Time {
	expires := time.Now().Add(expiresIn)
	return &expires
}
