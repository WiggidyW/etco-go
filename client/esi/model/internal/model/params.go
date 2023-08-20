package model

import "github.com/WiggidyW/weve-esi/client/esi/model/internal/naive"

type ModelParams[P naive.UrlParams, M any] struct {
	naive.NaiveParams[P]
	Model *M
}
