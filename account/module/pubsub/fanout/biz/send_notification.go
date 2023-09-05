package sendnotificationbiz

import (
	"context"

	"github.com/tronglv92/accounts/plugin/pubsub"
)

type sendNotificationBiz struct {
	ps pubsub.Pubsub
}

func NewSendNotificationBiz(ps pubsub.Pubsub) *sendNotificationBiz {
	return &sendNotificationBiz{

		ps: ps,
	}
}
func (biz *sendNotificationBiz) SendNotification(ctx context.Context) error {
	newMessage := pubsub.NewMessage(map[string]interface{}{

		"message": "notification",
	})
	// done := make(chan bool)
	_ = biz.ps.Publish(ctx, "fanout", "notification-exchange", "", "notification", newMessage)

	return nil
}
