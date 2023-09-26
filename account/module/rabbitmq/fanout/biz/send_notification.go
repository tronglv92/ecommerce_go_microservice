package sendnotificationbiz

import (
	"context"

	apprabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"
)

type sendNotificationBiz struct {
	ps apprabbitmq.Pubsub
}

func NewSendNotificationBiz(ps apprabbitmq.Pubsub) *sendNotificationBiz {
	return &sendNotificationBiz{

		ps: ps,
	}
}
func (biz *sendNotificationBiz) SendNotification(ctx context.Context) error {
	newMessage := apprabbitmq.NewMessage(map[string]interface{}{

		"message": "notification",
	})
	// done := make(chan bool)
	// _ = biz.ps.Publish(ctx, "fanout", "notification-exchange", "", "notification", newMessage)
	_ = biz.ps.PublishMsgToExchange(ctx, apprabbitmq.PublishConfig{
		ExchangeType: "fanout",
		ExchangeName: "notification-exchange",
		QueueName:    "",
		RoutingKey:   "notification",
		Data:         newMessage,
	})
	return nil
}
