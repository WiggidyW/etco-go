package items

import (
	"strings"

	"github.com/WiggidyW/weve-esi/error/esierror"
)

func RateLimited(err error) bool {
	statusErr, ok := err.(esierror.StatusError)
	if ok && statusErr.Code == LIMITED_CODE && strings.Contains(
		statusErr.EsiText,
		LIMITED_STR,
	) {
		return true
	}
	return false
}