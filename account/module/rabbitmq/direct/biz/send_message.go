package sendmessagebiz

import (
	"context"

	apprabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"
)

type sendMessageBiz struct {
	ps apprabbitmq.Pubsub
}

func NewSendMessageBiz(ps apprabbitmq.Pubsub) *sendMessageBiz {
	return &sendMessageBiz{

		ps: ps,
	}
}
func (biz *sendMessageBiz) SendMessage(ctx context.Context) error {
	newMessage := apprabbitmq.NewMessage(map[string]interface{}{

		"message": "hello restaurant",
	})
	// done := make(chan bool)
	// _ = biz.ps.Publish(ctx, "direct", "message-exchange", "message-queue", "message-key", newMessage)
	err := biz.ps.PublishMsgToExchange(ctx, apprabbitmq.PublishConfig{
		ExchangeType: "direct",
		ExchangeName: "message-exchange",
		QueueName:    "message-queue",
		RoutingKey:   "message-key",
		Data:         newMessage,
	})
	// err := biz.ps.PublishRetry(ctx, newMessage)
	return err
}
