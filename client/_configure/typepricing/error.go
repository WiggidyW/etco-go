package typepricing

import (
	"github.com/WiggidyW/etco-go/error/configerror"
)

func newError(
	errStr string,
) configerror.ErrPricingInvalid {
	return configerror.ErrPricingInvalid{ErrString: errStr}
}
