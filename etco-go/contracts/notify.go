package contracts

import (
	"context"
	"fmt"

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
	go func() {
		logger.MaybeErr(notifyNewContracts(
			x.Ctx(),
			newContracts.Buyback,
			build.CONTRACT_NOTIFICATIONS_BUYBACK_BASE_URL,
			"Buyback",
		))
	}()
	go func() {
		logger.MaybeErr(notifyNewContracts(
			x.Ctx(),
			newContracts.Shop,
			build.CONTRACT_NOTIFICATIONS_SHOP_BASE_URL,
			"Shop",
		))
	}()
	return nil
}

func notifyNewContracts(
	ctx context.Context,
	contracts map[string]Contract,
	baseUrl string,
	kindStr string,
) error {
	if len(contracts) == 0 {
		return nil
	}
	subject := "New " + kindStr + " Contracts"
	message := ""
	for code, contract := range contracts {
		if contract.Status == Outstanding {
			message += baseUrl + code + "\n"
		}
	}
	if message == "" {
		return nil
	}
	if err := notifier.ContractSend(
		ctx,
		subject,
		message,
	); err != nil {
		return fmt.Errorf(
			"failed to send 'New %s contracts' notification: %w",
			kindStr,
			err,
		)
	}
	return nil
}
