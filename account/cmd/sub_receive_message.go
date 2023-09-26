package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tronglv92/accounts/common"

	"github.com/tronglv92/accounts/plugin/opentelemetry"
	rabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"
	goservice "github.com/tronglv92/ecommerce_go_common"
)

type HasRestaurantId interface {
	GetRestaurantId() int
	GetUserId() int
}

var startSubReceiveMessageCmd = &cobra.Command{
	Use:   "sub-receive-message",
	Short: "Start a subscriber when user send message",
	Run: func(cmd *cobra.Command, args []string) {
		service := goservice.New(

			goservice.WithInitRunnable(rabbitmq.NewRabbitMQ(common.PluginRabbitMQ)),
			goservice.WithInitRunnable(opentelemetry.NewJaeger("ecommerce_recieve_message")),
		)

		if err := service.Init(); err != nil {
			log.Fatalln(err)
		}

		ps := service.MustGet(common.PluginRabbitMQ).(rabbitmq.Pubsub)

		// ctx := context.Background()

		ch, _ := ps.Consumer("direct", "message-exchange", "message-queue", "message-key")

		for msg := range ch {
			fmt.Printf("startSubReceiveMessageCmd msg: %v", msg)

		}
	},
}
