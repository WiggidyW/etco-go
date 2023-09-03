package structureinfo

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/etco-go/client/esi/model"
)

type StructureInfoParams struct {
	WebRefreshToken string
	StructureId     int64
}

type StructureInfoUrlParams struct {
	StructureId int64
}

func (p StructureInfoUrlParams) Url() string {
	return fmt.Sprintf(
		"%s/universe/structures/%d/?datasource=%s",
		model.BASE_URL,
		p.StructureId,
		model.DATASOURCE,
	)
}

func (StructureInfoUrlParams) Method() string {
	return http.MethodGet
}
