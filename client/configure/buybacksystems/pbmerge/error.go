package pbmerge

import (
	"fmt"

	"github.com/WiggidyW/eve-trading-co-go/error/configerror"
)

func newError(systemId int32, errStr string) configerror.ErrInvalid {
	return configerror.ErrInvalid{
		Err: configerror.ErrBuybackSystemInvalid{
			Err: fmt.Errorf(
				"'%d': %s",
				systemId,
				errStr,
			),
		},
	}
}
