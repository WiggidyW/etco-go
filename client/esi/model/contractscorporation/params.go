package contractscorporation

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/esi/model"
)

type ContractsCorporationParams struct {
	WebRefreshToken string
	CorporationId   int32
}

type ContractsCorporationUrlParams struct {
	CorporationId int32
}

func (p ContractsCorporationUrlParams) PageUrl(page *int) string {
	query := fmt.Sprintf(
		"%s/corporations/%d/contracts/?datasource=%s",
		model.BASE_URL,
		p.CorporationId,
		model.DATASOURCE,
	)
	query = model.AddQueryInt(query, "page", page)
	return query
}

func (ContractsCorporationUrlParams) Method() string {
	return http.MethodGet
}