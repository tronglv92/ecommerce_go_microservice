package sendemailbiz

import (
	"context"

	model "github.com/tronglv92/accounts/module/rabbitmq/topic/model"
	apprabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"
)

type sendEmailBiz struct {
	ps apprabbitmq.Pubsub
}

func NewSendEmailBiz(ps apprabbitmq.Pubsub) *sendEmailBiz {
	return &sendEmailBiz{

		ps: ps,
	}
}
func (biz *sendEmailBiz) SendEmail(ctx context.Context, data *model.MessageMail) error {
	newMessage := apprabbitmq.NewMessage(map[string]interface{}{

		"message": data.Message,
	})
	// done := make(chan bool)
	// error := biz.ps.Publish(ctx, "topic", "email-exchange", "", data.Routekey, newMessage)
	error := biz.ps.Publish(ctx, apprabbitmq.PublishConfig{
		ExchangeType: "topic",
		ExchangeName: "email-exchange",
		QueueName:    "",
		RoutingKey:   data.Routekey,
		Data:         newMessage,
	})
	return error
}
