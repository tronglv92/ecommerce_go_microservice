package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tronglv92/accounts/common"
	"github.com/tronglv92/accounts/plugin/kafka"

	goservice "github.com/tronglv92/ecommerce_go_common"
)

var startSubReceiveMessageFromKafkaCmd = &cobra.Command{
	Use:   "message-kafka",
	Short: "Start a subscriber when user send message",
	Run: func(cmd *cobra.Command, args []string) {
		service := goservice.New(

			goservice.WithInitRunnable(kafka.NewKafka(common.PluginKafka)),
		)

		if err := service.Init(); err != nil {
			log.Fatalln(err)
		}

		ps := service.MustGet(common.PluginKafka).(kafka.Pubsub)

		// topics, err := ps.GetTopics(context.Background(), "topic1231232")
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("startSubReceiveMessageFromKafkaCmd msg: %v \n", topics)
		// if len(topics) == 0 {
		// 	panic(errors.New("loi"))
		// }
		// ctx := context.Background()
		//Subscribe(ctx context.Context, config KafkaConsumerConfig) (ch <-chan *Message, close func(), errChan <-chan error)
		ch, _, errCh := ps.Subscribe(context.Background(), kafka.KafkaConsumerConfig{
			GroupId: "group1",
			Topic:   "topic1",
		})
		for {
			select {
			case msg := <-ch:
				fmt.Printf("startSubReceiveMessageFromKafkaCmd msg: %v \n", msg.Data)

			case err := <-errCh:
				fmt.Printf("startSubReceiveMessageFromKafkaCmd err: %v \n", err)

			}
		}

		// for msg := range ch {
		// 	fmt.Printf("startSubReceiveMessageFromKafkaCmd msg: %v \n", msg.Data)

		// }
	},
}
