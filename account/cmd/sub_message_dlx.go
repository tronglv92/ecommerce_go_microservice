package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tronglv92/accounts/common"

	rabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"
	goservice "github.com/tronglv92/ecommerce_go_common"
)

var startMessageDLXCmd = &cobra.Command{
	Use:   "sub-message-dlx",
	Short: "Start a subscriber when user send message",
	Run: func(cmd *cobra.Command, args []string) {
		service := goservice.New(

			goservice.WithInitRunnable(rabbitmq.NewRabbitMQ(common.PluginRabbitMQ)),
		)

		if err := service.Init(); err != nil {
			log.Fatalln(err)
		}

		ps := service.MustGet(common.PluginRabbitMQ).(rabbitmq.Pubsub)

		// ctx := context.Background()

		ch, _ := ps.Consumer("direct", "messageexdlx", "messagequeuedlx", "messagekeydlx")

		for msg := range ch {
			fmt.Printf("receive msg dlx: %v", msg)

		}
	},
}
