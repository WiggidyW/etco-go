package authorized

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
)

const (
	TOKEN_CHARACTER_MIN_EXPIRES time.Duration = 30 * time.Minute
	TOKEN_CHARACTER_BUF_CAP     int           = 0
)

func init() {
	keys.TypeStrTokenCharacter = cache.RegisterType[int32]("tokenCharacter", 0)
}

func GetTokenCharacter(
	x cache.Context,
	app esi.EsiApp,
	refreshToken string,
) (
	charId int32,
	expires time.Time,
	err error,
) {
	return tokenCharGet(x, app, refreshToken)
}
