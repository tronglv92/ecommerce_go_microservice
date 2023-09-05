package sendemailbiz

import (
	"context"

	model "github.com/tronglv92/accounts/module/pubsub/topic/model"
	"github.com/tronglv92/accounts/plugin/pubsub"
)

type sendEmailBiz struct {
	ps pubsub.Pubsub
}

func NewSendEmailBiz(ps pubsub.Pubsub) *sendEmailBiz {
	return &sendEmailBiz{

		ps: ps,
	}
}
func (biz *sendEmailBiz) SendEmail(ctx context.Context, data *model.MessageMail) error {
	newMessage := pubsub.NewMessage(map[string]interface{}{

		"message": data.Message,
	})
	// done := make(chan bool)
	error := biz.ps.Publish(ctx, "topic", "email-exchange", "", data.Routekey, newMessage)

	return error
}
