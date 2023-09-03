package contractitems

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/etco-go/client/esi/model"
)

type ContractItemsParams struct {
	WebRefreshToken string
	CorporationId   int32
	ContractId      int32
}

type ContractItemsUrlParams struct {
	CorporationId int32
	ContractId    int32
}

func (p ContractItemsUrlParams) Url() string {
	return fmt.Sprintf(
		"%s/corporations/%d/contracts/%d/items/?datasource=%s",
		model.BASE_URL,
		p.CorporationId,
		p.ContractId,
		model.DATASOURCE,
	)
}

func (ContractItemsUrlParams) Method() string {
	return http.MethodGet
}
