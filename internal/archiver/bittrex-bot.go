package archiver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func BittrexBot(botConfig *BotConfig) {

	defer func() {
		log.Printf("Orderbook Capture Complete: %v    ***\n", botConfig.Market)
	}()

	msg := fmt.Sprintf("Bot %s starting, Archive Orderbooks: %v, Archive Trades: %v\n",
		botConfig.Market,
		botConfig.ArchiveOrderbooks,
		botConfig.ArchiveTrades)
	fmt.Printf(msg)

	captureEndpointData := func(uri string) []byte {

		endPoint := "https://api.bittrex.com/v3/markets/" + botConfig.Market + uri

		resp, _ := http.Get(endPoint)
		failOnError(err, "Unable to access endpoint")
		defer resp.Body.Close()

		http_body, _ := io.ReadAll(resp.Body)

		// Need to publish both the bot's market and the received json body
		brokerMsg := BrokerMsg{
			Market:   botConfig.Market,
			HttpBody: http_body,
		}
		data, _ := json.Marshal(brokerMsg)

		return data
	}

	for {

		if botConfig.ArchiveOrderbooks {
			RabbitPublish("orderbooks", captureEndpointData("/orderbook?depth=25"))
		}

		if botConfig.ArchiveTrades {
			RabbitPublish("trades", captureEndpointData("/trades"))
		}

		// Take a rest
		time.Sleep(2 * time.Second)
	}
}
