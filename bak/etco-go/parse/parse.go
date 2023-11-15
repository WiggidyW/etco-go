package parse

import (
	"github.com/evepraisal/go-evepraisal/parsers"

	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/staticdb"
)

func Parse(text string) (
	knownItems []*proto.NamedBasicItem,
	unknownItems []*proto.NamedBasicItem,
) {
	return MapIntoKnownAndUnknown(ParseIntoMap(text))
}

func ParseIntoMap(text string) (parseMap map[string]int64) {
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

func MapIntoKnownAndUnknown(parseMap map[string]int64) (
	knownItems []*proto.NamedBasicItem,
	unknownItems []*proto.NamedBasicItem,
) {
	knownItems = make([]*proto.NamedBasicItem, 0, len(parseMap))
	unknownItems = make([]*proto.NamedBasicItem, 0, len(parseMap))

	var known bool
	for name, quantity := range parseMap {
		item := &proto.NamedBasicItem{
			Name:     name,
			Quantity: quantity,
		}
		item.TypeId, known = staticdb.NameToTypeId(name)
		if known {
			knownItems = append(knownItems, item)
		} else {
			unknownItems = append(unknownItems, item)
		}
	}

	return knownItems, unknownItems
}
