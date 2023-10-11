package ordersstructure

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/etco-go/client/esi/model"
)

type OrdersStructureParams struct {
	WebRefreshToken string
	StructureId     int64
}

type OrdersStructureUrlParams struct {
	StructureId int64
}

func (p OrdersStructureUrlParams) PageUrl(page *int) string {
	query := fmt.Sprintf(
		"%s/markets/structures/%d/?datasource=%s",
		model.BASE_URL,
		p.StructureId,
		model.DATASOURCE,
	)
	query = model.AddQueryInt(query, "page", page)
	return query
}

func (OrdersStructureUrlParams) Method() string {
	return http.MethodGet
}
