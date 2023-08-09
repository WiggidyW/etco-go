package tc

import "github.com/WiggidyW/weve-esi/staticdb"

type PricingInfo struct {
	market  *staticdb.Container[Market]
	pricing Pricing
}

func (p *PricingInfo) IsBuy() bool {
	return p.pricing.IsBuy
}

func (p *PricingInfo) Percentile() int {
	return int(p.pricing.Percentile)
}

func (p *PricingInfo) Modifier() uint8 {
	return p.pricing.Modifier
}

func (p *PricingInfo) MarketName() string {
	return p.getMarket().Name
}

func (p *PricingInfo) MarketRefreshToken() (string, bool) {
	market := p.getMarket()
	if market.RefreshToken == nil {
		return "", false
	}
	return *market.RefreshToken, true
}

func (p *PricingInfo) MarketLocationId() int64 {
	return p.getMarket().LocationId
}

func (p *PricingInfo) MarketIsStructure() bool {
	return p.getMarket().IsStructure
}

func newPricingInfo(p Pricing) *PricingInfo {
	return &PricingInfo{pricing: p}
}

func (p *PricingInfo) getMarket() Market {
	if p.market == nil {
		m := kVReaderMarket.UnsafeGet(p.pricing.MarketIndex)
		p.market = staticdb.NewContainer[Market](m)
	}
	return p.market.Inner
}
