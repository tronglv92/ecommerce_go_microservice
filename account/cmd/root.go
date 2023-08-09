package cmd

import (
	"fmt"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/tronglv92/accounts/common"
	"github.com/tronglv92/accounts/middleware"
	consul "github.com/tronglv92/accounts/plugin/consul"
	resty "github.com/tronglv92/accounts/plugin/resty"
	"github.com/tronglv92/ecommerce_go_common/plugin/storage/sdkgorm"

	handlers "github.com/tronglv92/accounts/cmd/handler"
	cardgrpcclient "github.com/tronglv92/accounts/plugin/grpc/card"
	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("food-delivery"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("mySql", common.DBMain)),
		goservice.WithInitRunnable(resty.NewRestService()),
		goservice.WithInitRunnable(cardgrpcclient.NewCardGrpcClient(common.PluginGrpcCardClient)),
		goservice.WithInitRunnable(consul.NewConsulClient(common.PluginConsul, "account")),
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

		// appContext := appctx.NewAppContext(db, s3Provider, secretKey, ps)
		service.HTTPServer().AddHandler(func(engine *gin.Engine) {

			engine.Use(middleware.Recover())
			engine.GET("/ping", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"data": "pong"})
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
