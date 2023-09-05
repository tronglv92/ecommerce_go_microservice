package cmd

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tronglv92/cards/common"
	"github.com/tronglv92/cards/middleware"
	"github.com/tronglv92/cards/plugin/storage/sdkgorm"

	"google.golang.org/grpc"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/nanmu42/gzip"
	handlers "github.com/tronglv92/cards/cmd/handler"
	cardGrpcbiz "github.com/tronglv92/cards/module/card/biz/grpc"
	"github.com/tronglv92/cards/plugin/consul"
	grpcService "github.com/tronglv92/cards/plugin/grpc"
	"github.com/tronglv92/cards/plugin/opentelemetry"
	cardgrpc "github.com/tronglv92/cards/proto/card"

	cardrepo "github.com/tronglv92/cards/module/card/repository"
	grpcstorage "github.com/tronglv92/cards/module/card/storage/grpc"
	goservice "github.com/tronglv92/ecommerce_go_common"
	"go.opentelemetry.io/otel/trace"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("food-delivery"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("mySql", common.DBMain)),
		goservice.WithInitRunnable(grpcService.NewGRPCServer(context.Background(), common.PluginGrpcServer)),
		goservice.WithInitRunnable(consul.NewConsulClient(common.PluginConsul, "card")),
		goservice.WithInitRunnable(opentelemetry.NewJaeger("ecommerce_go_card")),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start an food delivery service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		serviceLogger := service.Logger("service")

		grpServer := service.MustGet(common.PluginGrpcServer).(interface {
			SetRegisterHdl(hdl func(*grpc.Server))
			SetRegisterHdlGw(hdl func(context.Context, *runtime.ServeMux, *grpc.ClientConn))
		})
		grpServer.SetRegisterHdl(func(server *grpc.Server) {
			var dbConn interface{}

			for dbConn == nil {
				dbConn = service.MustGet(common.DBMain)
				fmt.Printf("dbConn %v \n", dbConn)
			}

			db, _ := dbConn.(sdkgorm.GormInterface)
			if db != nil {
				dbSession := db.Session()
				store := grpcstorage.NewSQLStore(db, dbSession)
				repo := cardrepo.NewListCardByCustomerIdRepo(store)
				biz := cardGrpcbiz.NewListCardByCustomerIdBiz(repo)

				cardgrpc.RegisterCardServiceServer(server, biz)
			}

			// user.RegisterDeviceTokenServiceServer(server, fcmGrpcServer.NewDeviceTokenGRPCBusiness(devicetokenstorage.NewSQLStore(dbConn)))

		})
		grpServer.SetRegisterHdlGw(func(ctx context.Context, gwmux *runtime.ServeMux, conn *grpc.ClientConn) {
			err := cardgrpc.RegisterCardServiceHandler(ctx, gwmux, conn)
			if err != nil {
				log.Fatalln("Failed to register gateway:", err)
			}
		})

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
