package sendmessagebiz

import (
	"context"

	apprabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"
)

type sendMessageBiz struct {
	ps apprabbitmq.Pubsub
}

func NewSendMessageDLXBiz(ps apprabbitmq.Pubsub) *sendMessageBiz {
	return &sendMessageBiz{

		ps: ps,
	}
}
func (biz *sendMessageBiz) SendMessageDLX(ctx context.Context) error {
	newMessage := apprabbitmq.NewMessage(map[string]interface{}{

		"message": "hello restaurant",
	})
	// done := make(chan bool)
	// _ = biz.ps.Publish(ctx, "direct", "message-exchange", "message-queue", "message-key", newMessage)
	err := biz.ps.PublishRetryWithDLX(ctx, apprabbitmq.PublishConfig{
		ExchangeType:    "direct",
		ExchangeName:    "message-retry-exchange",
		QueueName:       "message-retry-queue",
		RoutingKey:      "message-retry-key",
		ExchangeNameDLX: "messageexdlx",
		QueueNameDLX:    "messagequeuedlx",
		RoutingKeyDLX:   "messagekeydlx",
		Data:            newMessage,
	})
	return err
}
