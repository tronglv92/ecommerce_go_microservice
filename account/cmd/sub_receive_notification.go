package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tronglv92/accounts/common"

	rabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"
	goservice "github.com/tronglv92/ecommerce_go_common"
)

var startSubReceiveNotificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Start a subscriber when user send notification",
	Run: func(cmd *cobra.Command, args []string) {
		service := goservice.New(

			goservice.WithInitRunnable(rabbitmq.NewRabbitMQ(common.PluginRabbitMQ)),
		)

		if err := service.Init(); err != nil {
			log.Fatalln(err)
		}

		ps := service.MustGet(common.PluginRabbitMQ).(rabbitmq.Pubsub)

		// ctx := context.Background()

		ch, _ := ps.Subscribe("fanout", "notification-exchange", "", "notification")

		for msg := range ch {
			fmt.Printf("startSubReceiveNotificationCmd msg: %v", msg)

		}
	},
}
