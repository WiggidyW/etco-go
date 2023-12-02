package notifier

import (
	"context"
	"fmt"
	"sync"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/logger"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/discord"
)

const (
	CHARACTER_LIMIT int = 2000

	MAX_DISCORD_ATTEMPTS int           = 3
	DISCORD_SLEEP_INTVL  time.Duration = 5 * time.Second
)

var (
	mu *sync.RWMutex
)

func init() {
	mu = new(sync.RWMutex)
	mu.Lock()
	go func() {
		var wg sync.WaitGroup
		if build.DISCORD_BOT_TOKEN != "" && build.DISCORD_CHANNEL != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := discordInit(0)
				if err != nil {
					logger.Err(err.Error())
				}
			}()
		}
		wg.Wait()
		mu.Unlock()
	}()
}

func discordInit(attempt int) (err error) {
	discord := discord.New()
	discord.AddReceivers(build.DISCORD_CHANNEL)
	err = discord.AuthenticateWithBotToken(build.DISCORD_BOT_TOKEN)
	if err == nil {
		notify.UseServices(discord)
		return nil
	} else if attempt >= MAX_DISCORD_ATTEMPTS {
		return fmt.Errorf("failed to authenticate with discord: %w", err)
	}
	time.Sleep(DISCORD_SLEEP_INTVL)
	return discordInit(attempt + 1)
}

func urlListSend(
	subject string,
	ctx context.Context,
	baseUrl string,
	affixes ...string,
) (
	err error,
) {
	mu.RLock()
	defer mu.RUnlock()
	defer func() {
		if err != nil {
			err = fmt.Errorf(
				"failed to send '%s' notification: %w",
				subject,
				err,
			)
		}
	}()

	lenSubject := len(subject)
	if lenSubject > CHARACTER_LIMIT {
		err = fmt.Errorf("subject '%s' is too long", subject)
		return err
	}

	messages := make([]string, 0, 1)
	var message, affixMessage string
	var affixSumLen int

	for _, affix := range affixes {
		affixMessage = baseUrl + affix + "\n"
		affixSumLen = lenSubject + len(affixMessage)
		if len(message)+affixSumLen > CHARACTER_LIMIT {
			if affixSumLen > CHARACTER_LIMIT {
				err = fmt.Errorf("affix '%s' is too long to send", affix)
				return err
			} else {
				messages = append(messages, message)
				message = affixMessage
			}
		} else {
			message += affixMessage
		}
	}
	messages = append(messages, message)

	for _, message := range messages {
		err := notify.Send(ctx, subject, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func BuybackContractsSend(
	ctx context.Context,
	affixes ...string,
) error {
	return urlListSend(
		"New Buyback Contracts",
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
		ctx,
		build.SHOP_CONTRACT_NOTIFICATIONS_BASE_URL,
		affixes...,
	)
}

func HaulContractsSend(
	ctx context.Context,
	affixes ...string,
) error {
	return urlListSend(
		"New Haul Contracts",
		ctx,
		build.HAUL_CONTRACT_NOTIFICATIONS_BASE_URL,
		affixes...,
	)
}

func PurchasesSend(
	ctx context.Context,
	affixes ...string,
) error {
	return urlListSend(
		"New Purchases",
		ctx,
		build.PURCHASE_NOTIFICATIONS_BASE_URL,
		affixes...,
	)
}
