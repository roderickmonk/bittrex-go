package archiver

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	// "os"
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

			resp, err := http.Get(
				"https://api.bittrex.com/v3/markets/" +
					botConfig.Market +
					"/orderbook?depth=25")
			if err != nil {
				log.Panic(err)
			}
			defer resp.Body.Close()

			http_body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("io.ReadAll: ", err)
				panic(err)
			}

			brokerMsg := BrokerMsg{
				Market:   botConfig.Market,
				HttpBody: http_body,
			}
			bytes, _ := json.Marshal(brokerMsg)

			RabbitPublish("orderbooks", bytes)
		}

		if botConfig.ArchiveTrades {

			resp, err := http.Get("https://api.bittrex.com/v3/markets/" + botConfig.Market + "/trades")
			if err != nil {
				log.Panic(err)
			}
			defer resp.Body.Close()

			http_body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("io.ReadAll: ", err)
				panic(err)
			}

			brokerMsg := BrokerMsg{
				Market:   botConfig.Market,
				HttpBody: http_body,
			}
			bytes, _ := json.Marshal(brokerMsg)

			RabbitPublish("trades", bytes)
		}

		time.Sleep(sleepTime)
	}
}
