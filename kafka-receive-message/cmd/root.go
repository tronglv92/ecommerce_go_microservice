package cmd

import (
	"context"
	"fmt"

	"log"
	"os"
	"time"

	"github.com/tronglv92/kafka-receive-message/common"

	"github.com/spf13/cobra"
	 "github.com/tronglv92/ecommerce_go_common"
	kafka "github.com/tronglv92/kafka-receive-message/plugin/kafka"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("food-delivery"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(kafka.NewKafka(common.PluginKafka)),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start an food delivery service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		serviceLogger := service.Logger("service")

		initServiceWithRetry(service, 10)

		ps := service.MustGet(common.PluginKafka).(kafka.Pubsub)

		// ctx := context.Background()

		ch, _ := ps.Subscribe(context.Background(), "group1", "topic1")

		for msg := range ch {
			fmt.Printf("startSubReceiveMessageFromKafkaCmd msg: %v \n", msg.Data)

		}

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}

	},
}

func Execute() {
	// TransAddPoint outenv as a sub command
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initServiceWithRetry(s goservice.Service, retry int) {
	var err error
	for i := 1; i <= retry; i++ {
		if err = s.Init(); err != nil {
			s.Logger("service").Errorf("error when starting service: %s", err.Error())
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}
