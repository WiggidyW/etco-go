package contracts

import (
	"context"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/notifier"
	"github.com/WiggidyW/etco-go/remotedb"
)

func getAndNotifyNewContracts(
	x cache.Context,
	contracts Contracts,
) (
	err error,
) {
	var newContracts remotedb.NewContracts[Contract, Contract]
	newContracts, _, err = remotedb.GetNewContracts(
		x,
		contracts.BuybackContracts,
		contracts.ShopContracts,
	)
	if err != nil {
		return err
	}
	if build.BUYBACK_CONTRACT_NOTIFICATIONS {
		go func() {
			logger.MaybeErr(notifyNewContracts(
				x.Ctx(),
				newContracts.Buyback,
				notifier.BuybackContractsSend,
			))
		}()
	}
	if build.SHOP_CONTRACT_NOTIFICATIONS {
		go func() {
			logger.MaybeErr(notifyNewContracts(
				x.Ctx(),
				newContracts.Shop,
				notifier.ShopContractsSend,
			))
		}()
	}
	return nil
}

func notifyNewContracts(
	ctx context.Context,
	contracts map[string]Contract,
	send func(context.Context, ...string) error,
) error {
	outstandingCodes := newOutstandingContractCodes(contracts)
	if len(outstandingCodes) > 0 {
		return send(ctx, outstandingCodes...)
	} else {
		return nil
	}
}

func newOutstandingContractCodes(
	contracts map[string]Contract,
) (
	outstandingCodes []string,
) {
	outstandingCodes = make([]string, 0, len(contracts))
	for code, contract := range contracts {
		if contract.Status == Outstanding {
			outstandingCodes = append(outstandingCodes, code)
		}
	}
	return outstandingCodes
}
