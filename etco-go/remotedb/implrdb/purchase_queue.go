package implrdb

type RawPurchaseQueue = map[int64][]string

type CodeAndLocationId struct {
	Code       string
	LocationId int64
}

func (c CodeAndLocationId) GetCode() string      { return c.Code }
func (c CodeAndLocationId) GetLocationId() int64 { return c.LocationId }
