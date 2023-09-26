package apprabbitmq

import (
	"context"
)

type Pubsub interface {
	PublishMsgToExchange(ctx context.Context, publishConfig PublishConfig) error
	// PublishRetryWithDLX(ctx context.Context, data *Message) error

	Consumer(exchangeType string, exchangeName string, queueName string, routingKey string) (ch <-chan *Message, close func())
	PublishMsgToQueue(ctx context.Context, config QueueDelayExpireConfig) error
	// Consumer(exchangeType string, exchangeName string, queueName string, routingKey string, exchangeNameDLX string, routingKeyDLX string) (ch <-chan *Message, close func())
	// PublishRetry(ctx context.Context, data *Message) error
	//UnSubcribe(ctx context.Context, channel Channel) error
}
