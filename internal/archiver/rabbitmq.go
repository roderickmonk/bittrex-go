package archiver

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var connection *amqp.Connection
var channel *amqp.Channel
var q amqp.Queue

func RabbitConnect() {

	connection, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = connection.Channel()
	failOnError(err, "Failed to open a RabbitMQ channel")

}

func RabbitPublish(routing_key string, data []byte) {

	err = channel.Publish(
		"GENERAL", // exchange
		routing_key,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	failOnError(err, "Failed to publish a message")
}

func RabbitClose() {
	channel.Close()
	connection.Close()
}
