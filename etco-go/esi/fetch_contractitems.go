package esi

import (
	"errors"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/error/esierror"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
)

const (
	CONTRACT_ITEMS_ENTRIES_NUM_RETRIES int = ESI_NUM_RETRIES
)

var (
	contractItemsEntriesRetry *fetch.Retry = &fetch.Retry{
		Retries:     CONTRACT_ITEMS_ENTRIES_NUM_RETRIES,
		ShouldRetry: contractItemsEntriesShouldRetry,
	}
)

func contractItemsEntriesShouldRetry(err error) bool {
	var statusErr esierror.StatusError
	if errors.As(err, &statusErr) {
		if rateLimited(statusErr) {
			return true
		} else if esiShouldRetryInner(statusErr) {
			return true
		}
	}
	return false
}

func contractItemsEntriesGet(
	x cache.Context,
	contractId int32,
) (
	rep []ContractItemsEntry,
	expires time.Time,
	err error,
) {
	url := contractItemsEntriesUrl(contractId)
	return fetch.HandleFetch(
		x,
		nil,
		contractItemsEntriesGetFetchFunc(contractId, url),
		contractItemsEntriesRetry,
	)
}

func contractItemsEntriesGetNewRep() *[]ContractItemsEntry {
	rep := make([]ContractItemsEntry, 0, CONTRACT_ITEMS_ENTRIES_MAKE_CAP)
	return &rep
}

func contractItemsEntriesGetFetchFunc(
	contractId int32,
	url string,
) fetch.Fetch[[]ContractItemsEntry] {
	return func(x cache.Context) (
		rep []ContractItemsEntry,
		expires time.Time,
		_ *postfetch.Params,
		err error,
	) {
		ciRateLimiterStart()
		rep, expires, err = getModel(
			x,
			url,
			CONTRACT_ITEMS_ENTRIES_METHOD,
			EsiAuthCorp,
			contractItemsEntriesGetNewRep,
		)
		go ciRateLimiterDone()
		if err != nil {
			var statusErr esierror.StatusError
			if errors.As(err, &statusErr) && statusErr.Code == 404 {
				err = nil
			} else {
				return nil, expires, nil, err
			}
		}
		expires = fetch.CalcExpires(
			expires,
			CONTRACT_ITEMS_ENTRIES_MIN_EXPIRES_IN,
		)
		return rep, expires, nil, err
	}
}
