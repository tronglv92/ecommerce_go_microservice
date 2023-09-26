package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tronglv92/accounts/common"

	"github.com/tronglv92/accounts/plugin/opentelemetry"
	rabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"
	goservice "github.com/tronglv92/ecommerce_go_common"
)

var startRetryMessageCmd = &cobra.Command{
	Use:   "sub-retry-message",
	Short: "Start a retry message",
	Run: func(cmd *cobra.Command, args []string) {
		service := goservice.New(

			goservice.WithInitRunnable(rabbitmq.NewRabbitMQ(common.PluginRabbitMQ)),
			goservice.WithInitRunnable(opentelemetry.NewJaeger("sub retry rabbit mq")),
		)

		if err := service.Init(); err != nil {
			log.Fatalln(err)
		}

		rabbit := service.MustGet(common.PluginRabbitMQ).(rabbitmq.Pubsub)

		ctx := context.Background()

		ch, _ := rabbit.Consumer("direct", "ex_dlx", "q_dlx", "route_dlx")

		for msg := range ch {
			fmt.Printf("receive msg.RetryCount: %v", msg.RetryCount)
			msg.RetryCount++
			if msg.RetryCount == 1 {

				rabbit.PublishMsgToQueue(ctx, rabbitmq.QueueDelayExpireConfig{
					Message: msg,

					TTL:         10000,
					QueueName:   "q_delay_ex_1",
					ExchangeDLX: "ex_mcrv",
					RouteDLX:    "route_mcrv",
				})
			} else if msg.RetryCount == 2 {
				rabbit.PublishMsgToQueue(ctx, rabbitmq.QueueDelayExpireConfig{
					Message: msg,

					TTL:         15000,
					QueueName:   "q_delay_ex_2",
					ExchangeDLX: "ex_mcrv",
					RouteDLX:    "route_mcrv",
				})
			} else if msg.RetryCount == 3 {
				rabbit.PublishMsgToQueue(ctx, rabbitmq.QueueDelayExpireConfig{
					Message: msg,

					TTL:         20000,
					QueueName:   "q_delay_ex_3",
					ExchangeDLX: "ex_mcrv",
					RouteDLX:    "route_mcrv",
				})
			} else {

			}
		}
	},
}
