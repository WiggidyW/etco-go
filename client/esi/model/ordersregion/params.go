package ordersregion

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/esi/model"
)

type OrdersRegionParams struct {
	RegionId int32
	TypeId   int32
	IsBuy    bool
}

type OrdersRegionUrlParams struct {
	RegionId  int32
	TypeId    *int32
	OrderType *string
}

func (p OrdersRegionUrlParams) PageUrl(page *int) string {
	query := fmt.Sprintf(
		"%s/markets/%d/orders/?datasource=%s",
		model.BASE_URL,
		p.RegionId,
		model.DATASOURCE,
	)
	query = model.AddQueryInt32(query, "type_id", p.TypeId)
	query = model.AddQueryString(query, "order_type", p.OrderType)
	query = model.AddQueryInt(query, "page", page)
	return query
}

func (OrdersRegionUrlParams) Method() string {
	return http.MethodGet
}
