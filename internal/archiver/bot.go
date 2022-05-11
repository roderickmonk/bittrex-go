package archiver

import (
	"encoding/json"
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

	for {

		if botConfig.ArchiveOrderbooks {

			// msg := fmt.Sprintf("Orderbooks: %v", botConfig.Market)
			// RabbitPublish("orderbooks", msg)

			orderbook := Orderbook{
				Ts:     time.Now(),
				Market: botConfig.Market,
				Buy:    []Entry{{1.0, 5.1}, {2.0, 5.2}},
				Sell:   []Entry{{3.0, 5.3}, {4.0, 5.4}},
			}

			// fmt.Println("orderbook:", orderbook)

			json, err := json.Marshal(orderbook)
			if err != nil {
				panic(err)
			}

			// fmt.Println("marshalled", string(json))
			// fmt.Println("marshalled (raw)", json)

			RabbitPublish("orderbooks", json)
		}

		if botConfig.ArchiveTrades {

			trade := Trade{
				Ts:         time.Now(),
				Market:     botConfig.Market,
				Id:         fmt.Sprintf("%v", rand.Int()),
				ExecutedAt: time.Now(),
				R:          1.0,
				Q:          5.0,
				TakerSide:  "buy",
			}

			json, err := json.Marshal(trade)
			if err != nil {
				panic(err)
			}

			RabbitPublish("trades", json)
		}

		time.Sleep(sleepTime)
	}
}
