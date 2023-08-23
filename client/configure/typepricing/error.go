package typepricing

import (
	"github.com/WiggidyW/weve-esi/error/configerror"
)

func newError(
	errStr string,
) configerror.ErrPricingInvalid {
	return configerror.ErrPricingInvalid{ErrString: errStr}
}
