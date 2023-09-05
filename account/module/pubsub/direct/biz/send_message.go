package sendmessagebiz

import (
	"context"

	"github.com/tronglv92/accounts/plugin/pubsub"
)

type sendMessageBiz struct {
	ps pubsub.Pubsub
}

func NewSendMessageBiz(ps pubsub.Pubsub) *sendMessageBiz {
	return &sendMessageBiz{

		ps: ps,
	}
}
func (biz *sendMessageBiz) SendMessage(ctx context.Context) error {
	newMessage := pubsub.NewMessage(map[string]interface{}{

		"message": "hello restaurant",
	})
	// done := make(chan bool)
	_ = biz.ps.Publish(ctx, "direct", "message-exchange", "message-queue", "message-key", newMessage)

	return nil
}
