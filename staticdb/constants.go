package staticdb

const (
	CORPORATION_ID    int32  = 0 // TODO
	WEB_REFRESH_TOKEN string = "TODO"
	BUYBACK_VERSION   string = "TODO"
	SHOP_VERSION      string = "TODO"
)

type HashSet[T comparable] map[T]struct{}

func (hs HashSet[T]) Has(k T) bool {
	_, ok := hs[k]
	return ok
}
