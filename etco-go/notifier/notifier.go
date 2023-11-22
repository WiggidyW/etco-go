package notifier

import (
	"context"

	"github.com/nikoksr/notify"
)

var (
	ContractNotifier notify.Notifier
)

func init() {
	ContractNotifier = notify.New()
}

func ContractSend(
	ctx context.Context,
	subject, message string,
) error {
	return ContractNotifier.Send(ctx, subject, message)
}
