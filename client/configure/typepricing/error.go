package typepricing

import (
	"github.com/WiggidyW/eve-trading-co-go/error/configerror"
)

func newError(
	errStr string,
) configerror.ErrPricingInvalid {
	return configerror.ErrPricingInvalid{ErrString: errStr}
}
