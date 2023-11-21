package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/error/esierror"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/proto/protoerr"
)

const (
	TOKEN_CHARACTER_MIN_EXPIRES time.Duration = 30 * time.Minute
	TOKEN_CHARACTER_BUF_CAP     int           = 0
)

func init() {
	keys.TypeStrTokenCharacter = cache.RegisterType[int32]("tokenCharacter", 0)
}

func ProtoGetTokenCharacter(
	x cache.Context,
	app esi.EsiApp,
	refreshToken string,
) (
	charId int32,
	expires time.Time,
	err error,
) {
	charId, expires, err = tokenCharGet(x, app, refreshToken)
	if err != nil {
		var statusErr esierror.StatusError
		if errors.As(err, &statusErr) &&
			statusErr.Code == 400 &&
			strings.Contains(statusErr.EsiText, "invalid_grant") {
			err = protoerr.New(protoerr.TOKEN_INVALID, err)
		}
	}
	return charId, expires, err
}
