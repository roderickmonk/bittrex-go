package archiver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func BrokerReceiver(mongoClient *mongo.Client, routing_key string) {
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

	collection := MongoClient.Database("history").Collection(routing_key)

	archive := func(http_body []byte) {

		switch routing_key {

		case "orderbooks":

			var brokerMsg BrokerMsg

			if err := json.Unmarshal(http_body, &brokerMsg); err != nil {
				fmt.Println("Unmarshal: ", err)
				panic(err)
			}
			fmt.Println("market: ", brokerMsg.Market)

			fmt.Println("Orderbook Received")
			var orderbook BittrexOrderbook
			if err := json.Unmarshal(brokerMsg.HttpBody, &orderbook); err != nil {
				fmt.Println("Unmarshal: ", err)
				panic(err)
			}

			var bids []MongoOrderbookEntry
			for _, e := range orderbook.Bid {
				r, _ := strconv.ParseFloat(e.Rate, 64)
				q, _ := strconv.ParseFloat(e.Quantity, 64)

				m := MongoOrderbookEntry{r, q}
				bids = append(bids, m)
			}
			fmt.Println("bids: ", bids)

			var asks []MongoOrderbookEntry
			for _, e := range orderbook.Ask {
				r, _ := strconv.ParseFloat(e.Rate, 64)
				q, _ := strconv.ParseFloat(e.Quantity, 64)

				m := MongoOrderbookEntry{r, q}
				asks = append(asks, m)
			}
			fmt.Println("asks: ", asks)

			save_orderbook := MongoOrderbook{
				Ts:     time.Now(),
				Market: brokerMsg.Market,
				Bid:    bids,
				Ask:    asks,
			}
			_, err = collection.InsertOne(context.TODO(), save_orderbook)
			if err != nil {
				panic(err)
			}

		case "trades":

			var brokerMsg BrokerMsg

			if err := json.Unmarshal(http_body, &brokerMsg); err != nil {
				fmt.Println("Unmarshal: ", err)
				panic(err)
			}

			fmt.Println("market: ", brokerMsg.Market)

			var trades []BittrexTrade
			if err := json.Unmarshal(brokerMsg.HttpBody, &trades); err != nil {
				fmt.Println("Unmarshal: ", err)
				panic(err)
			}

			for i, t := range trades {
				fmt.Println(i, t.ExecutedAt, t.Id, t.R, t.Q, t.TakerSide)
			}

			var mongoTrades []MongoTrade
			for i, t := range trades {
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
				fmt.Println(i, t.ExecutedAt, t.Id, t.R, t.Q, t.TakerSide)
			}
			fmt.Println("mongoTrades:\n", mongoTrades)

			var saveTrades = MongoTrades{
				Ts:     time.Now(),
				Market: brokerMsg.Market,
				Trades: mongoTrades,
			}
			_, err = collection.InsertOne(context.TODO(), saveTrades)
			if err != nil {
				panic(err)
			}
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
