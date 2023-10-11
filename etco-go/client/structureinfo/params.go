package structureinfo

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
)

type StructureInfoParams struct {
	StructureId int64
}

func (p StructureInfoParams) CacheKey() string {
	return cachekeys.StructureInfoCacheKey(p.StructureId)
}
