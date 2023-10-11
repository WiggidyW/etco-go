package webtocore

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_ATTRS_TIMEOUT = 60 * time.Second
)

type WebAttrs struct {
	// hex version of the CRC32C
	CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER string
	CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER  string
	CHECKSUM_WEB_MARKETS                          string
	CHECKSUM_WEB_SHOP_LOCATIONS                   string
	CHECKSUM_WEB_BUYBACK_SYSTEMS                  string

	VERSION_TIME   time.Time
	VERSION_STRING string // time string
}

func TransceiveDownloadWebAttrs(
	ctx context.Context,
	bucketClient *b.BucketClient,
	chnSend chanresult.ChanSendResult[WebAttrs],
) error {
	webAttrs, err := DownloadWebAttrs(ctx, bucketClient)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(webAttrs)
	}
}

func DownloadWebAttrs(
	ctx context.Context,
	bucketClient *b.BucketClient,
) (WebAttrs, error) {
	ctx, cancel := context.WithTimeout(ctx, WEB_ATTRS_TIMEOUT)
	defer cancel()

	webAttrs := &WebAttrs{}

	chnAttrs := chanresult.NewChanResult[AttrsWithVariant](ctx, 5, 0)
	chnSendAttrs, chnRecvAttrs := chnAttrs.Split()

	go transceiveAttrs(
		ctx,
		ATTR_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
		bucketClient.ReadAttrsWebBuybackSystemTypeMapsBuilder,
		chnSendAttrs,
	)
	go transceiveAttrs(
		ctx,
		ATTR_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
		bucketClient.ReadAttrsWebShopLocationTypeMapsBuilder,
		chnSendAttrs,
	)
	go transceiveAttrs(
		ctx,
		ATTR_WEB_BUYBACK_SYSTEMS,
		bucketClient.ReadAttrsWebBuybackSystems,
		chnSendAttrs,
	)
	go transceiveAttrs(
		ctx,
		ATTR_WEB_SHOP_LOCATIONS,
		bucketClient.ReadAttrsWebShopLocations,
		chnSendAttrs,
	)
	go transceiveAttrs(
		ctx,
		ATTR_WEB_MARKETS,
		bucketClient.ReadAttrsWebMarkets,
		chnSendAttrs,
	)

	for i := 0; i < 5; i++ {
		attrs, err := chnRecvAttrs.Recv()
		if err != nil {
			return WebAttrs{}, err
		} else {
			webAttrs.addAttrs(attrs)
		}
	}

	return *webAttrs, nil
}

func (wa *WebAttrs) addAttrs(attrs AttrsWithVariant) {
	if attrs.Attrs == nil { // no object exists
		return
	}

	switch attrs.AttrVariant {
	case ATTR_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER:
		wa.CHECKSUM_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER =
			attrsChecksum(attrs.Attrs)
	case ATTR_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER:
		wa.CHECKSUM_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER =
			attrsChecksum(attrs.Attrs)
	case ATTR_WEB_BUYBACK_SYSTEMS:
		wa.CHECKSUM_WEB_BUYBACK_SYSTEMS = attrsChecksum(attrs.Attrs)
	case ATTR_WEB_SHOP_LOCATIONS:
		wa.CHECKSUM_WEB_SHOP_LOCATIONS = attrsChecksum(attrs.Attrs)
	case ATTR_WEB_MARKETS:
		wa.CHECKSUM_WEB_MARKETS = attrsChecksum(attrs.Attrs)
	}

	// version is just the latest created or updated time
	var versionTime time.Time
	if attrs.Attrs.Updated.After(attrs.Attrs.Created) {
		versionTime = attrs.Attrs.Updated
	} else {
		versionTime = attrs.Attrs.Created
	}
	if versionTime.After(wa.VERSION_TIME) {
		wa.VERSION_TIME = attrs.Attrs.Created
		wa.VERSION_STRING = timeToVersion(versionTime)
	}
}

type AttrsWithVariant struct {
	AttrVariant AttrVariant
	Attrs       *b.Attrs
}

type AttrVariant uint8

const (
	ATTR_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER AttrVariant = iota
	ATTR_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER
	ATTR_WEB_BUYBACK_SYSTEMS
	ATTR_WEB_SHOP_LOCATIONS
	ATTR_WEB_MARKETS
)

func attrsChecksum(attrs *b.Attrs) string {
	return fmt.Sprintf("%x", attrs.CRC32C)
}

func timeToVersion(time time.Time) string {
	return time.UTC().Format("06-01-02_15:04:05_UTC")
}

func transceiveAttrs(
	ctx context.Context,
	attrVariant AttrVariant,
	fetch func(ctx context.Context) (*b.Attrs, error),
	chnSendAttrs chanresult.ChanSendResult[AttrsWithVariant],
) error {
	if attrs, err := fetch(ctx); err != nil {
		return chnSendAttrs.SendErr(err)
	} else {
		return chnSendAttrs.SendOk(AttrsWithVariant{
			AttrVariant: attrVariant,
			Attrs:       attrs,
		})
	}
}
