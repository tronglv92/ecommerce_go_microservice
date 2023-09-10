package apprabbitmq

import (
	"context"
)

type Pubsub interface {
	Publish(ctx context.Context, publishConfig PublishConfig) error
	PublishRetryWithDLX(ctx context.Context, publishConfig PublishConfig) error

	Subscribe(exchangeType string, exchangeName string, queueName string, routingKey string) (ch <-chan *Message, close func())
	Consumer(exchangeType string, exchangeName string, queueName string, routingKey string, exchangeNameDLX string, routingKeyDLX string) (ch <-chan *Message, close func())
	//UnSubcribe(ctx context.Context, channel Channel) error
}
