package protoutil

import (
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/proto"
)

func NewPBCfgMarkets(
	webMarkets map[b.MarketName]b.WebMarket,
) (
	pbMarkets map[string]*proto.CfgMarket,
) {
	pbMarkets = make(
		map[string]*proto.CfgMarket,
		len(webMarkets),
	)
	for marketName, webMarket := range webMarkets {
		pbMarkets[marketName] = NewPBCfgMarket(webMarket)
	}
	return pbMarkets
}

func NewPBCfgMarket(
	webMarket b.WebMarket,
) (
	pbMarket *proto.CfgMarket,
) {
	if webMarket.RefreshToken != nil {
		return &proto.CfgMarket{
			RefreshToken: *webMarket.RefreshToken,
			LocationId:   webMarket.LocationId,
			IsStructure:  webMarket.IsStructure,
		}
	} else {
		return &proto.CfgMarket{
			// RefreshToken: "",
			LocationId:  webMarket.LocationId,
			IsStructure: webMarket.IsStructure,
		}
	}
}
