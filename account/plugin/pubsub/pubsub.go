package pubsub

import "context"

type Pubsub interface {
	Publish(ctx context.Context, exchangeType string, exchangeName string, queueName string, routingKey string, data *Message) error

	Subscribe(exchangeType string, exchangeName string, queueName string, routingKey string) (ch <-chan *Message, close func())
	//UnSubcribe(ctx context.Context, channel Channel) error
}
