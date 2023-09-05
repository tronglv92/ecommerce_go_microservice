package kafka

import "context"

type Pubsub interface {
	Publish(ctx context.Context, topic string, messages ...*Message) error

	Subscribe(ctx context.Context, config KafkaConsumerConfig) (ch <-chan *Message, close func(), errChan <-chan error)

	CreateTopics(ctx context.Context, topics ...string) error
	CheckTopicsCreated(ctx context.Context, topicsName ...string) error
}
