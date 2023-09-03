package structureinfo

import "github.com/WiggidyW/etco-go/error/esierror"

type StructureInfo struct {
	Forbidden bool
	Name      string // "" if not authorized
	SystemId  int32  // -1 if not authorized
}

func Forbidden(err error) bool {
	statusErr, ok := err.(esierror.StatusError)
	if ok && statusErr.Code == 403 {
		return true
	}
	return false
}
