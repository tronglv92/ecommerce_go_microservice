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
	err := biz.ps.PublishMsgToExchange(ctx, apprabbitmq.PublishConfig{
		ExchangeType: "direct",
		ExchangeName: "ex_mcrv",
		QueueName:    "q_mcsv",
		RoutingKey:   "route_mcrv",
		Data:         newMessage,

		ExchangeNameDLX: "ex_dlx",
		RoutingKeyDLX:   "route_dlx",
		TTL:             10000,
	})
	return err
}
