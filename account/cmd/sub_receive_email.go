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

var startSubReceiveEmailCmd = &cobra.Command{
	Use:   "email",
	Short: "Start a subscriber when user send notification",
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

		for _, key := range args {
			ch, _ := ps.Consumer("topic", "email-exchange", "", key)

			for msg := range ch {
				fmt.Printf("startSubReceiveEmailCmd msg: %v", msg)

			}
		}

	},
}
