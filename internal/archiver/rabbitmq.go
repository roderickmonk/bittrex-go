package archiver

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var connection *amqp.Connection
var err error
var channel *amqp.Channel
var q amqp.Queue

func RabbitConnect() {

	connection, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = connection.Channel()
	failOnError(err, "Failed to open a channel")

}

func RabbitPublish(routing_key string, body string) {

	err = channel.Publish(
		"GENERAL", // exchange
		routing_key,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf("Sent to Broker: %s\n", body)
}

func RabbitClose() {
	channel.Close()
	connection.Close()
}
