package kafka

import "context"

type Pubsub interface {
	Publish(ctx context.Context, topic string, messages ...*Message) error

	Subscribe(ctx context.Context, groupId string, topic string) (ch <-chan *Message, close func())
	//UnSubcribe(ctx context.Context, channel Channel) error
}
