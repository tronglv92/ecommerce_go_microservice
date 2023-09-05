package appgrpc

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tronglv92/cards/common"
	"github.com/tronglv92/ecommerce_go_common/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	prefix        string
	port          int
	server        *grpc.Server
	logger        logger.Logger
	ctx           context.Context
	registerHdl   func(*grpc.Server)
	registerHdlGw func(context.Context, *runtime.ServeMux, *grpc.ClientConn)
}

func NewGRPCServer(ctx context.Context, prefix string) *grpcServer {
	return &grpcServer{prefix: prefix, ctx: ctx}
}
func (s *grpcServer) SetRegisterHdl(hdl func(*grpc.Server)) {
	s.registerHdl = hdl
}
func (s *grpcServer) SetRegisterHdlGw(hdl func(context.Context, *runtime.ServeMux, *grpc.ClientConn)) {
	s.registerHdlGw = hdl
}
func (s *grpcServer) GetPrefix() string {
	return s.prefix
}
func (s *grpcServer) Get() interface{} {
	return s
}

func (s *grpcServer) Name() string {
	return s.prefix
}
func (s *grpcServer) InitFlags() {
	flag.IntVar(&s.port, s.GetPrefix()+"-port", 50051, "Port of gRPC service")
}

func (s *grpcServer) Configure() error {
	s.logger = logger.GetCurrent().GetLogger(s.prefix)
	s.logger.Infoln("Setup gRPC service:", s.prefix)
	s.logger.Infoln("Setup gRPC service:", s.port)

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	reflection.Register(s.server)
	return nil
}
func (s *grpcServer) Run() error {

	time.Sleep(time.Second * 2)
	_ = s.Configure()
	if s.registerHdl != nil {
		s.logger.Infoln("registering services...")
		s.registerHdl(s.server)
	}

	s.logger.Infoln("registering services success %v", s.registerHdl)
	address := fmt.Sprintf("0.0.0.0:%d", s.port)
	lis, err := net.Listen("tcp", address)
	// if address != '' {
	// 	s.logger.Info("Connected gRPC service at ", address)
	// }

	if err != nil {
		s.logger.Errorln("Error %v", err)
		// return err
	}
	go func() {
		defer common.AppRecover()
		s.server.Serve(lis)

	}()

	conn, err := grpc.DialContext(
		s.ctx,
		address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	if s.registerHdlGw != nil {
		s.logger.Infoln("registering services gateway...")
		s.registerHdlGw(s.ctx, gwmux, conn)
	}

	gwServer := &http.Server{
		Addr:    ":2000",
		Handler: gwmux,
	}
	go func() {
		defer common.AppRecover()
		gwServer.ListenAndServe()

	}()

	return nil
}
func (s *grpcServer) Stop() <-chan bool {
	c := make(chan bool)

	go func() {
		s.server.Stop()
		c <- true
		s.logger.Infoln("Stopped")
	}()
	return c
}
