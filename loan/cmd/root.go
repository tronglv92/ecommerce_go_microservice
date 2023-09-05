package cmd

import (
	"fmt"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/nanmu42/gzip"
	"github.com/tronglv92/loans/common"
	"github.com/tronglv92/loans/middleware"
	"github.com/tronglv92/loans/plugin/consul"
	"github.com/tronglv92/loans/plugin/opentelemetry"
	"github.com/tronglv92/loans/plugin/storage/sdkgorm"
	"go.opentelemetry.io/otel/trace"

	goservice "github.com/tronglv92/ecommerce_go_common"
	handlers "github.com/tronglv92/loans/cmd/handler"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("food-delivery"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("mySql", common.DBMain)),
		goservice.WithInitRunnable(consul.NewConsulClient(common.PluginConsul, "loan")),
		goservice.WithInitRunnable(opentelemetry.NewJaeger("ecommerce_go_loan")),
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

		tracer := service.MustGet(common.PluginOpenTelemetry).(trace.Tracer)
		db := service.MustGet(common.DBMain).(sdkgorm.GormInterface)
		err := db.RegisterGormCallbacks(tracer)
		if err != nil {
			panic(err)
		}

		// appContext := appctx.NewAppContext(db, s3Provider, secretKey, ps)
		service.HTTPServer().AddHandler(func(engine *gin.Engine) {

			engine.Use(middleware.Recover())
			engine.Use(gzip.DefaultHandler().Gin)
			engine.GET("/ping", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"data": "pong"})
			})
			handlers.MainRoute(engine, service)
			handlers.InternalRoute(engine, service)
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
