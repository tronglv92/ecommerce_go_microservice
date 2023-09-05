package kafkabiz

import (
	"context"

	"github.com/tronglv92/accounts/plugin/kafka"
)

type sendMessageBiz struct {
	ps kafka.Pubsub
}

func NewSendMessageBiz(ps kafka.Pubsub) *sendMessageBiz {
	return &sendMessageBiz{

		ps: ps,
	}
}
func (biz *sendMessageBiz) SendMessage(ctx context.Context) error {
	newMessage := kafka.NewMessage(map[string]interface{}{

		"message": "hello kafka",
	})
	// done := make(chan bool)
	error := biz.ps.Publish(ctx, "topic1", newMessage)

	return error
}
