package buyback

import (
	"context"
	"hash"
	"hash/fnv"
	"time"

	"github.com/WiggidyW/weve-esi/client/appraisal"
	bm "github.com/WiggidyW/weve-esi/client/market/buyback"
	wbanon "github.com/WiggidyW/weve-esi/client/remotedb/appraisal/writebuyback/anonymous"
	wbchar "github.com/WiggidyW/weve-esi/client/remotedb/appraisal/writebuyback/character"
	"github.com/WiggidyW/weve-esi/staticdb"
	"github.com/WiggidyW/weve-esi/util"
)

type BuybackAppraisalClient struct {
	writeCharClient wbchar.SAC_WriteBuybackCharacterAppraisalClient
	writeAnonClient wbanon.WriteBuybackAnonAppraisalClient
	marketClient bm.BuybackPriceClient
}

func (bac BuybackAppraisalClient) Fetch(
	ctx context.Context,
	params BuybackAppraisalParams,
) (*appraisal.BuybackAppraisal, error) {
	systemInfoPtr := staticdb.GetBuybackSystemInfo(params.SystemId)
	if systemInfoPtr == nil {
		return NewRejectedAppraisal(params), nil
	}
	systemInfo := *systemInfoPtr

	// fetch the appraisal items in a goroutine
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := util.NewChanResult[*appraisal.BuybackParentItem](
		ctx,
	).Split()
	for _, item := range params.Items {
		go bac.fetchOne(ctx, item, systemInfo, chnSend)
	}

	// create a hash code if we're saving the appraisal
	var hasher hash.Hash64 = nil
	if params.Save {
		hasher = fnv.New64a()
	}

	// initialize the appraisal
	appraisal := &appraisal.BuybackAppraisal{
		Items: make([]appraisal.BuybackParentItem, len(params.Items)),
		Price: 0.0,
		Version: staticdb.BUYBACK_VERSION,
		SystemId: params.SystemId,
		CharacterId: params.CharacterId,
		

	// collect the results
}

func (bac BuybackAppraisalClient) fetchOne(
	ctx context.Context,
	item appraisal.BasicItem,
	systemInfo staticdb.BuybackSystemInfo,
	chnSend util.ChanSendResult[*appraisal.BuybackParentItem],
) error {
	if apprItem, err := bac.marketClient.Fetch(
		ctx,
		bm.BuybackPriceParams{
			BuybackSystemInfo: systemInfo,
			TypeId:            item.TypeId,
			Quantity: 	item.Quantity,
		},
	); err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(apprItem)
	}
}

// bm "github.com/WiggidyW/weve-esi/client/market/buyback"
