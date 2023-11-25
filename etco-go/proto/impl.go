package proto

type Nullable interface {
	IsNil() bool
}

func (a *BuybackAppraisal) IsNil() bool { return a == nil }
func (a *ShopAppraisal) IsNil() bool    { return a == nil }

type Appraisal interface {
	Nullable
	GetCharacterId() int32
	ClearCharacterId()
}

func (a *BuybackAppraisal) ClearCharacterId() { a.CharacterId = 0 }
func (a *ShopAppraisal) ClearCharacterId()    { a.CharacterId = 0 }
