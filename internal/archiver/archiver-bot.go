package archiver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func ArchiveBot(mongoClient *mongo.Client, routing_key string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"GENERAL", // name
		"direct",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "GENERAL", routing_key)
	err = ch.QueueBind(
		q.Name,      // queue name
		routing_key, // routing key
		"GENERAL",   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	rabbitMsgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	collection := mongoClient.Database("history").Collection(routing_key)

	level1Parse := func(msg []byte) (market string, httpBody *[]byte) {

		var brokerMsg BrokerMsg

		err = json.Unmarshal(msg, &brokerMsg)
		failOnError(err, "Unmarshal Error")

		return brokerMsg.Market, &brokerMsg.HttpBody
	}

	archive := func(msg []byte) {

		switch routing_key {

		case "orderbooks":

			fmt.Println("Orderbook Received")

			// Nested unmarshalling
			market, httpBody := level1Parse(msg)

			var orderbook BittrexOrderbook
			err = json.Unmarshal(*httpBody, &orderbook)
			failOnError(err, "Unmarshal Error")

			// Massage the Bittrex orderbook data ready to be archived to Mongo
			var bids []MongoOrderbookEntry
			for _, e := range orderbook.Bid {
				r, _ := strconv.ParseFloat(e.Rate, 64)
				q, _ := strconv.ParseFloat(e.Quantity, 64)

				m := MongoOrderbookEntry{r, q}
				bids = append(bids, m)
			}
			// fmt.Println("bids: ", bids)

			var asks []MongoOrderbookEntry
			for _, e := range orderbook.Ask {
				r, _ := strconv.ParseFloat(e.Rate, 64)
				q, _ := strconv.ParseFloat(e.Quantity, 64)

				m := MongoOrderbookEntry{r, q}
				asks = append(asks, m)
			}
			// fmt.Println("asks: ", asks)

			document := MongoOrderbook{
				Ts:     time.Now(),
				Market: market,
				Bid:    bids,
				Ask:    asks,
			}
			_, err = collection.InsertOne(context.TODO(), document)
			failOnError(err, "Unable to insert document into orderbooks collection")

		case "trades":

			fmt.Println("Trades Received")

			// Nested unmarshalling
			market, httpBody := level1Parse(msg)

			var trades []BittrexTrade
			err := json.Unmarshal(*httpBody, &trades)
			failOnError(err, "Unmarshal Error")

			// Massage the Bittrex trades data ready to be archived to Mongo
			var mongoTrades []MongoTrade
			for _, t := range trades {
				executedAt, _ := time.Parse(time.RFC3339, t.ExecutedAt)
				r, _ := strconv.ParseFloat(t.R, 64)
				q, _ := strconv.ParseFloat(t.Q, 64)
				mongoTrades = append(mongoTrades, MongoTrade{
					Id:         t.Id,
					ExecutedAt: executedAt,
					R:          r,
					Q:          q,
					TakerSide:  t.TakerSide,
				})
			}

			var document = MongoTrades{
				Ts:     time.Now(),
				Market: market,
				Trades: mongoTrades,
			}
			_, err = collection.InsertOne(context.TODO(), document)
			failOnError(err, "Unable to insert document into trades collection")

		}
	}

	go func() {
		for msg := range rabbitMsgs {
			archive(msg.Body)
		}
	}()

	log.Printf("Waiting for %s", routing_key)

	var forever chan struct{}
	<-forever
}
