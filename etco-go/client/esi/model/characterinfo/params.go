package characterinfo

import (
	"fmt"
	"net/http"

	"github.com/WiggidyW/etco-go/client/cachekeys"
	"github.com/WiggidyW/etco-go/client/esi/model"
)

type CharacterInfoParams struct {
	CharacterId int32
}

func (p CharacterInfoParams) CacheKey() string {
	return cachekeys.CharacterInfoCacheKey(p.CharacterId)
}

type CharacterInfoUrlParams struct {
	CharacterId int32
}

func (p CharacterInfoUrlParams) Url() string {
	return fmt.Sprintf(
		"%s/characters/%d/?datasource=%s",
		model.BASE_URL,
		p.CharacterId,
		model.DATASOURCE,
	)
}

func (CharacterInfoUrlParams) Method() string {
	return http.MethodGet
}
