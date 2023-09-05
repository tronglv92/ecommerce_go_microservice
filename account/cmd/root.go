package cmd

import (
	"context"
	"errors"
	"fmt"

	"log"
	"os"
	"time"

	"github.com/tronglv92/accounts/common"
	"github.com/tronglv92/accounts/middleware"
	consul "github.com/tronglv92/accounts/plugin/consul"
	"github.com/tronglv92/accounts/plugin/opentelemetry"
	resty "github.com/tronglv92/accounts/plugin/resty"
	"github.com/tronglv92/accounts/plugin/storage/sdkgorm"
	"github.com/tronglv92/accounts/plugin/storage/sdkmgo"
	"github.com/tronglv92/accounts/plugin/storage/sdkredis"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
	"github.com/nanmu42/gzip"
	"github.com/spf13/cobra"
	handlers "github.com/tronglv92/accounts/cmd/handler"
	cardgrpcclient "github.com/tronglv92/accounts/plugin/grpc/card"
	kafka "github.com/tronglv92/accounts/plugin/kafka"
	rabbitmq "github.com/tronglv92/accounts/plugin/pubsub/rabbitmq"
	goservice "github.com/tronglv92/ecommerce_go_common"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("food-delivery"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(kafka.NewKafka(common.PluginKafka)),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("mySql", common.DBMain)),
		goservice.WithInitRunnable(sdkmgo.NewMongoDB("mongoDB", common.DBMongo)),
		goservice.WithInitRunnable(resty.NewRestService()),
		goservice.WithInitRunnable(cardgrpcclient.NewCardGrpcClient(common.PluginGrpcCardClient)),
		goservice.WithInitRunnable(consul.NewConsulClient(common.PluginConsul, "account")),
		goservice.WithInitRunnable(opentelemetry.NewJaeger("ecommerce_go_account")),
		goservice.WithInitRunnable(sdkredis.NewRedisDB("redis", common.PluginRedis)),
		goservice.WithInitRunnable(rabbitmq.NewRabbitMQ(common.PluginRabbitMQ)),
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

		// register gorm callback to track
		tracer := service.MustGet(common.PluginOpenTelemetry).(trace.Tracer)
		db := service.MustGet(common.DBMain).(sdkgorm.GormInterface)

		err := db.RegisterGormCallbacks(tracer)
		if err != nil {
			panic(err)
		}

		ps := service.MustGet(common.PluginKafka).(kafka.Pubsub)
		ps.CreateTopics(context.Background(), "topic2")

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {

			engine.Use(middleware.Recover())
			engine.Use(gzip.DefaultHandler().Gin)
			// engine.Use(otelgin.Middleware("my-server"))
			// engine.Use(middleware.OpenTelemetryMiddleware("my-server"))

			engine.GET("/ping", func(ctx *gin.Context) {
				// ctx.JSON(http.StatusBadRequest, gin.H{"error": "pong"})
				panic(errors.New("has error"))
			})
			engine.POST("/ping", func(ctx *gin.Context) {
				// ctx.JSON(http.StatusBadRequest, gin.H{"error": "pong"})
				panic(errors.New("has error"))
			})
			handlers.MainRoute(engine, service)
			engine.StaticFile("/demo/", "./demo.html")

		})

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}

	},
}

func Execute() {
	// TransAddPoint outenv as a sub command
	rootCmd.AddCommand(outEnvCmd)
	rootCmd.AddCommand(startSubReceiveMessageCmd)
	rootCmd.AddCommand(startSubReceiveNotificationCmd)
	rootCmd.AddCommand(startSubReceiveEmailCmd)
	rootCmd.AddCommand(startSubReceiveMessageFromKafkaCmd)
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
