package notifier

import (
	"context"
	"fmt"

	build "github.com/WiggidyW/etco-go/buildconstants"

	"github.com/nikoksr/notify"
)

var (
	BuybackContractNotifier notify.Notifier
	ShopContractNotifier    notify.Notifier
	PurchaseNotifier        notify.Notifier
)

func init() {
	BuybackContractNotifier = notify.New()
	ShopContractNotifier = notify.New()
	PurchaseNotifier = notify.New()
}

func urlListSend(
	subject string,
	notifier notify.Notifier,
	ctx context.Context,
	baseUrl string,
	affixes ...string,
) error {
	message := ""
	for _, affix := range affixes {
		message += baseUrl + affix + "\n"
	}
	err := notifier.Send(ctx, subject, message)
	if err != nil {
		return fmt.Errorf("failed to send '%s' notification: %w", subject, err)
	} else {
		return nil
	}
}

func BuybackContractsSend(
	ctx context.Context,
	affixes ...string,
) error {
	return urlListSend(
		"New Buyback Contracts",
		BuybackContractNotifier,
		ctx,
		build.BUYBACK_CONTRACT_NOTIFICATIONS_BASE_URL,
		affixes...,
	)
}

func ShopContractsSend(
	ctx context.Context,
	affixes ...string,
) error {
	return urlListSend(
		"New Shop Contracts",
		ShopContractNotifier,
		ctx,
		build.SHOP_CONTRACT_NOTIFICATIONS_BASE_URL,
		affixes...,
	)
}

func PurchasesSend(
	ctx context.Context,
	affixes ...string,
) error {
	return urlListSend(
		"New Purchases",
		PurchaseNotifier,
		ctx,
		build.PURCHASE_NOTIFICATIONS_BASE_URL,
		affixes...,
	)
}
