package streaming

import (
	"context"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/entries"
	"github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/head"
	pe "github.com/WiggidyW/eve-trading-co-go/client/esi/model/internal/pageentries"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type StreamingPageEntriesClient[P pe.UrlPageParams, E any] struct {
	entriesClient entries.EntriesClient[staticUrlParams, E]
	headClient    head.HeadClient[staticUrlParams]
}

// func (spec StreamingPageEntriesClient[P, E]) Fetch(
// 	ctx context.Context,
// 	params pe.NaivePageParams[P],
// ) util.ChanRecvResult[HeadRepWithChan[cache.ExpirableData[[]E]]] {
// 	// create receiving and sending channels
// 	chnSend, chnRecv := util.
// 		NewChanResult[HeadRepWithChan[cache.ExpirableData[[]E]]](
// 		ctx,
// 	).Split()

// 	// call fetchStream() in a separate goroutine
// 	go func() {
// 		// send the HeadRepWithChan to the receiving channel
// 		hrwc, err := spec.fetchStream(ctx, params)
// 		if err != nil {
// 			chnSend.SendErr(err)
// 		} else {
// 			chnSend.SendOk(*hrwc)
// 		}
// 	}()

// 	return chnRecv
// }

func (spec StreamingPageEntriesClient[P, E]) Fetch(
	ctx context.Context,
	params pe.NaivePageParams[P],
) (*HeadRepWithChan[E], error) {
	// first, do a head fetch to get num pages
	headRep, err := spec.fetchHead(ctx, params)
	if err != nil {
		return nil, err
	}

	// then, create the page receiver channel
	chnRecv := spec.fetchPages(ctx, params, headRep.Data())

	return &HeadRepWithChan[E]{
		NumPages: headRep.Data(),
		Expires:  headRep.Expires(),
		ChanRecv: chnRecv,
	}, nil
}

// fetches the number of pages
func (spec StreamingPageEntriesClient[P, E]) fetchHead(
	ctx context.Context,
	params pe.NaivePageParams[P],
) (*cache.ExpirableData[int], error) {
	return spec.headClient.Fetch(
		ctx,
		newNaivePageParams[P](params, nil),
	)
}

// fetches a single page
func (spec StreamingPageEntriesClient[P, E]) fetchPage(
	ctx context.Context,
	params pe.NaivePageParams[P],
	page *int,
) (*cache.ExpirableData[[]E], error) {
	return spec.entriesClient.Fetch(
		ctx,
		newNaivePageParams[P](params, page),
	)
}

// returns a channel that will receive all page fetch results
func (spec StreamingPageEntriesClient[P, E]) fetchPages(
	ctx context.Context,
	params pe.NaivePageParams[P],
	pages int,
) util.ChanRecvResult[cache.ExpirableData[[]E]] {
	// create receiving and sending channels
	chnSend, chnRecv := util.NewChanResult[cache.ExpirableData[[]E]](
		ctx,
	).Split()

	// fetch each page in a separate goroutine
	for page := 1; page <= pages; page++ {
		go func(page int) {
			rep, err := spec.fetchPage(ctx, params, &page)
			if err != nil {
				chnSend.SendErr(err)
			} else {
				chnSend.SendOk(*rep)
			}
		}(page)
	}

	return chnRecv
}

// // returns a HeadRepWithChan
// func (spec StreamingPageEntriesClient[P, E]) fetchStream(
// 	ctx context.Context,
// 	params pe.NaivePageParams[P],
// ) (*HeadRepWithChan[cache.ExpirableData[[]E]], error) {
// }
