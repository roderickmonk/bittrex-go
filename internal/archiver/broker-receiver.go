package archiver

import (
	// "encoding/json"
	"fmt"
	"log"
	// "time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
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

	archive := func(json []byte) {

		var doc interface{}
		err := bson.UnmarshalExtJSON(json, true, &doc)
		if err != nil {
			// handle error
		}

		res, insertErr := collection.InsertOne(MongoCtx, doc)
		if insertErr != nil {
			log.Fatal(insertErr)
		}
		fmt.Println(res)
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
