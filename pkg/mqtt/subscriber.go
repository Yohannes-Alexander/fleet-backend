package mqtt

import (
	"context"
	"encoding/json"
	"log"

	paho "github.com/eclipse/paho.mqtt.golang"
	"fleet-backend/internal/dto"
	"fleet-backend/internal/usecase"
)

type Subscriber struct{
	client paho.Client
	uc     *usecase.FleetUsecase
	topic  string
}

func NewSubscriber(brokerURL, clientID, topic string, uc *usecase.FleetUsecase) *Subscriber {
	opts := paho.NewClientOptions().AddBroker(brokerURL).SetClientID(clientID)
	opts.SetAutoReconnect(true)
	client := paho.NewClient(opts)
	return &Subscriber{client: client, uc: uc, topic: topic}
}

func (s *Subscriber) Start(ctx context.Context) error {
	if token := s.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Printf("[MQTT] connected, subscribing %s", s.topic)
	if token := s.client.Subscribe(s.topic, 1, s.handle); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	go func() {
		<-ctx.Done()
		s.client.Disconnect(250)
	}()
	return nil
}

func (s *Subscriber) handle(_ paho.Client, msg paho.Message) {
	var m dto.MqttLocationMessage
	if err := json.Unmarshal(msg.Payload(), &m); err != nil {
		log.Printf("[MQTT] invalid payload: %v", err); return
	}
	if err := s.uc.IngestLocation(context.Background(), m); err != nil {
		log.Printf("[UC] ingest error: %v", err)
	}
}
