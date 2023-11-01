package raw_

import "github.com/WiggidyW/etco-go/error/esierror"

func shouldRetry(err error) bool {
	statusErr, ok := err.(esierror.StatusError)
	return ok && (statusErr.Code == 502 || statusErr.Code == 504)
}
