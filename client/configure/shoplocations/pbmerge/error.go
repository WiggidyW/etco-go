package pbmerge

import (
	"fmt"

	"github.com/WiggidyW/eve-trading-co-go/error/configerror"
)

func newError(locationId int64, errStr string) configerror.ErrInvalid {
	return configerror.ErrInvalid{
		Err: configerror.ErrShopLocationInvalid{
			Err: fmt.Errorf(
				"'%d': %s",
				locationId,
				errStr,
			),
		},
	}
}