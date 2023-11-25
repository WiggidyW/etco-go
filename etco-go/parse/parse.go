package parse

import (
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"

	"github.com/evepraisal/go-evepraisal/parsers"
)

type ParseMap map[string]int64

func Parse(text string) (parseMap ParseMap) {
	allParserResult, _ := parsers.AllParser(parsers.StringToInput(text))
	parserResults := allParserResult.(*parsers.MultiParserResult).Results
	parseMap = make(map[string]int64, len(parserResults))

	for _, sub_result := range parserResults {
		switch r := sub_result.(type) {
		case *parsers.AssetList:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.CargoScan:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.Contract:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.EFT:
			parseMap[r.Ship] += 1
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.Fitting:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.Industry:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.Listing:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.LootHistory:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.PI:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.SurveyScan:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.ViewContents:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.MiningLedger:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.MoonLedger:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.HeuristicResult:
			for _, item := range r.Items {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.DScan:
			for _, item := range r.Items {
				parseMap[item.Name] += 1
			}
		case *parsers.Compare:
			for _, item := range r.Items {
				parseMap[item.Name] += 1
			}
		case *parsers.Wallet:
			for _, item := range r.ItemizedTransactions {
				parseMap[item.Name] += item.Quantity
			}
		case *parsers.Killmail:
			for _, item := range r.Dropped {
				parseMap[item.Name] += item.Quantity
			}
			for _, item := range r.Destroyed {
				parseMap[item.Name] += item.Quantity
			}
		}
	}

	return parseMap
}

func (pm ParseMap) ToProto(
	r *protoregistry.ProtoRegistry,
) (
	knownItems []*proto.NamedBasicItem,
	unknownItems []*proto.NamedBasicItem,
) {
	knownItems = make([]*proto.NamedBasicItem, 0, len(pm))
	unknownItems = make([]*proto.NamedBasicItem, 0, len(pm))

	var named *proto.NamedTypeId
	var exists bool
	var item *proto.NamedBasicItem
	for name, quantity := range pm {
		named, exists = r.AddTypeByName(name)
		item = &proto.NamedBasicItem{TypeId: named, Quantity: quantity}
		if exists {
			knownItems = append(knownItems, item)
		} else {
			unknownItems = append(unknownItems, item)
		}
	}

	return knownItems, unknownItems
}

type ProtoParseRep struct {
	KnownItems   []*proto.NamedBasicItem
	UnknownItems []*proto.NamedBasicItem
}

func ProtoParse(
	// _ cache.Context, // (uncomment + import cache if needed for implementing interface)
	r *protoregistry.ProtoRegistry,
	text string,
) (
	rep ProtoParseRep,
	// err error, // (uncomment if needed for implementing interface)
) {
	parseMap := Parse(text)
	rep.KnownItems, rep.UnknownItems = parseMap.ToProto(r)
	return rep
}
