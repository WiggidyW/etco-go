package market

type Price interface {
	Price() float64
	Desc() string
}
