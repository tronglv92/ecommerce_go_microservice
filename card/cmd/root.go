package cmd

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tronglv92/cards/common"
	"github.com/tronglv92/cards/middleware"
	"github.com/tronglv92/ecommerce_go_common/plugin/storage/sdkgorm"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	handlers "github.com/tronglv92/cards/cmd/handler"
	cardGrpcbiz "github.com/tronglv92/cards/module/card/biz/grpc"
	grpcService "github.com/tronglv92/cards/plugin/grpc"
	cardgrpc "github.com/tronglv92/cards/proto/card"

	cardrepo "github.com/tronglv92/cards/module/card/repository"
	cardstorage "github.com/tronglv92/cards/module/card/storage/gorm"
	goservice "github.com/tronglv92/ecommerce_go_common"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("food-delivery"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("mySql", common.DBMain)),
		goservice.WithInitRunnable(grpcService.NewGRPCServer(context.Background(), common.PluginGrpcServer)),
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
			}
			service.Logger("service").Errorf("error dbConn: %s", dbConn)

			db, ok := dbConn.(*gorm.DB)
			if ok {
				store := cardstorage.NewSQLStore(db)
				repo := cardrepo.NewListCardByCustomerIdRepo(store)
				biz := cardGrpcbiz.NewListCardByCustomerIdBiz(repo)

				cardgrpc.RegisterCardServiceServer(server, biz)
			}

			// user.RegisterDeviceTokenServiceServer(server, fcmGrpcServer.NewDeviceTokenGRPCBusiness(devicetokenstorage.NewSQLStore(dbConn)))

		})
		grpServer.SetRegisterHdlGw(func(ctx context.Context, gwmux *runtime.ServeMux, conn *grpc.ClientConn) {
			err := cardgrpc.RegisterCardServiceHandler(context.Background(), gwmux, conn)
			if err != nil {
				log.Fatalln("Failed to register gateway:", err)
			}
		})

		initServiceWithRetry(service, 10)

		// appContext := appctx.NewAppContext(db, s3Provider, secretKey, ps)
		service.HTTPServer().AddHandler(func(engine *gin.Engine) {

			engine.Use(middleware.Recover())
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
