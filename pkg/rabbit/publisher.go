package rabbit

import (
	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct{
	ch *amqp091.Channel
}

func NewPublisher(conn *amqp091.Connection, exchange, kind, queue, routingKey string) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil { return nil, err }
	if err := ch.ExchangeDeclare(exchange, kind, true, false, false, false, nil); err != nil { return nil, err }
	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil { return nil, err }
	if err := ch.QueueBind(queue, routingKey, exchange, false, nil); err != nil { return nil, err }
	return &Publisher{ch: ch}, nil
}

func (p *Publisher) Publish(exchange, routingKey string, body []byte) error {
	return p.ch.Publish(exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
