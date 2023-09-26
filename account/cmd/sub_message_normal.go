package cmd

import (
	"github.com/spf13/cobra"
)

var startMessageNormalCmd = &cobra.Command{
	Use:   "sub-message-normal",
	Short: "Start a subscriber when user send message",
	Run: func(cmd *cobra.Command, args []string) {
		// service := goservice.New(

		// 	goservice.WithInitRunnable(rabbitmq.NewRabbitMQ(common.PluginRabbitMQ)),
		// )

		// if err := service.Init(); err != nil {
		// 	log.Fatalln(err)
		// }

		// ps := service.MustGet(common.PluginRabbitMQ).(rabbitmq.Pubsub)

		// // ctx := context.Background()

		// ch, _ := ps.Consumer("direct", "message-retry-exchange", "message-retry-queue", "message-retry-key", "messageexdlx", "messagekeydlx")

		// for msg := range ch {
		// 	fmt.Printf("receive msg normal: %v", msg)

		// }
	},
}
