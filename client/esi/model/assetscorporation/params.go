package assetscorporation

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/eve-trading-co-go/client/esi/model"
)

type AssetsCorporationParams struct {
	WebRefreshToken string
	CorporationId   int32
}

type AssetsCorporationUrlParams struct {
	CorporationId int32
}

func (p AssetsCorporationUrlParams) PageUrl(page *int) string {
	query := fmt.Sprintf(
		"%s/corporations/%d/assets/?datasource=%s",
		model.BASE_URL,
		p.CorporationId,
		model.DATASOURCE,
	)
	query = model.AddQueryInt(query, "page", page)
	return query
}

func (AssetsCorporationUrlParams) Method() string {
	return http.MethodGet
}
