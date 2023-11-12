package staticdb

import (
	"fmt"

	kvreader_ "github.com/WiggidyW/etco-go/staticdb/kvreaders_"
)

func init() {
	go loadAllCrashOnError()
}

func loadAllCrashOnError() {
	chnErr := make(chan error, 16)
	go initSendErr(kvreader_.InitKVReaderNameToTypeId, chnErr)
	go initSendErr(kvreader_.InitKVReaderBuybackSystemTypeMaps, chnErr)
	go initSendErr(kvreader_.InitKVReaderShopLocationTypeMaps, chnErr)
	go initSendErr(kvreader_.InitKVReaderTypeDataMap, chnErr)
	go initSendErr(kvreader_.InitKVReaderBuybackSystems, chnErr)
	go initSendErr(kvreader_.InitKVReaderShopLocations, chnErr)
	go initSendErr(kvreader_.InitKVReaderBannedFlagSets, chnErr)
	go initSendErr(kvreader_.InitKVReaderMarkets, chnErr)
	go initSendErr(kvreader_.InitKVReaderPricings, chnErr)
	go initSendErr(kvreader_.InitKVReaderCategories, chnErr)
	go initSendErr(kvreader_.InitKVReaderGroups, chnErr)
	go initSendErr(kvreader_.InitKVReaderMarketGroups, chnErr)
	go initSendErr(kvreader_.InitKVReaderRegions, chnErr)
	go initSendErr(kvreader_.InitKVReaderSystems, chnErr)
	go initSendErr(kvreader_.InitKVReaderStations, chnErr)
	go initSendErr(kvreader_.InitKVReaderTypeVolumes, chnErr)
	for i := 0; i < 16; i++ {
		if err := <-chnErr; err != nil {
			panic(fmt.Errorf("error loading static data: %w", err).Error())
		}
	}
}

func initSendErr(
	init func() error,
	chn chan<- error,
) {
	chn <- init()
}
