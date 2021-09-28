package main

import (
	"log"
	"os"
	"strings"

	connector "github.com/streadway/amqp"
)

func main() {
	conn, err := connector.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("%s: %s", "Could not connect RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Could not open channel", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare Queue", err)
	}

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		connector.Publishing{
			DeliveryMode: connector.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		log.Fatalf("%s: %s", "Could not publish messages", err)
	}
	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
