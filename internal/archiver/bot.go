package archiver

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(99))

func Bot(botConfig *BotConfig) {

	deferred := func() {
		log.Printf("Orderbook Capture Complete: %v    ***\n", botConfig.Market)
	}

	defer deferred()

	msg := fmt.Sprintf("Bot %s starting, Archive Orderbooks: %v, Archive Trades: %v\n", botConfig.Market, botConfig.ArchiveOrderbooks, botConfig.ArchiveTrades)
	fmt.Printf(msg)

	sleepTime := time.Duration(r.Int31n(1000)) * 5_000_000
	// fmt.Println("timeBetween: ", sleepTime)

	for i := 0; i < 5; i++ {

		if botConfig.ArchiveOrderbooks {
			msg := fmt.Sprintf("Orderbooks: %v", botConfig.Market)
			RabbitPublish("orderbooks", msg)
		}

		if botConfig.ArchiveTrades {
			msg := fmt.Sprintf("Trades: %v", botConfig.Market)
			RabbitPublish("trades", msg)
		}

		time.Sleep(sleepTime)
	}
}
