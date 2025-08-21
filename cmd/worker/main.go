package main

import (
	"log"

	"fleet-backend/config"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.Load()
	conn, err := amqp091.Dial(cfg.RabbitURL)
	if err != nil { log.Fatal(err) }
	ch, err := conn.Channel()
	if err != nil { log.Fatal(err) }

	_, _ = ch.QueueDeclare(cfg.RabbitQueue, true, false, false, false, nil)

	msgs, err := ch.Consume(cfg.RabbitQueue, "", true, false, false, false, nil)
	if err != nil { log.Fatal(err) }

	log.Printf("[WORKER] listening queue=%s", cfg.RabbitQueue)
	for msg := range msgs {
		log.Printf("[WORKER] geofence alert: %s", string(msg.Body))
		// TODO: push ke service lain, email, dsb.
	}
}
