package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/felipedias-dev/fullcycle-go-expert-clean-architecture/pkg/events"
	"github.com/streadway/amqp"
)

type OrderCreateHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrderCreateHandler(rabbitMQChannel *amqp.Channel) *OrderCreateHandler {
	return &OrderCreateHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrderCreateHandler) HandleEvent(event events.EventInterface, wq *sync.WaitGroup) error {
	defer wq.Done()
	fmt.Printf("Order created: %v", event.GetPayload())
	jsoOutput, err := json.Marshal(event.GetPayload())
	if err != nil {
		return err
	}

	msgRabbitMq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsoOutput,
	}

	err = h.RabbitMQChannel.Publish(
		"amqp.direct",
		"",
		false,
		false,
		msgRabbitMq,
	)
	if err != nil {
		return err
	}
	return nil
}
