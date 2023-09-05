package appgrpc

import (
	"flag"

	cardgrpc "github.com/tronglv92/accounts/proto/card"
	"github.com/tronglv92/ecommerce_go_common/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	logger      logger.Logger
	prefix      string
	url         string
	gwSupported bool
	gwPort      int
	client      cardgrpc.CardServiceClient
	cc          *grpc.ClientConn
}

func NewCardGrpcClient(prefix string) *grpcClient {
	return &grpcClient{
		prefix: prefix,
	}
}

func (client *grpcClient) GetPrefix() string {
	return client.prefix
}

func (client *grpcClient) Get() interface{} {
	return client.client
}

func (client *grpcClient) Name() string {
	return client.prefix
}

func (client *grpcClient) InitFlags() {
	flag.StringVar(&client.url, client.GetPrefix()+"-url", "localhost:50051", "URL connect to grpc server")
}

func (client *grpcClient) Configure() error {
	client.logger = logger.GetCurrent().GetLogger(client.Name())
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	cc, err := grpc.Dial(client.url, opts,
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	if err != nil {
		return err
	}
	client.logger.Infoln("grpc client connected success")
	client.cc = cc
	client.client = cardgrpc.NewCardServiceClient(cc)

	return nil
}

func (client *grpcClient) Run() error {
	return client.Configure()
}

func (client *grpcClient) Stop() <-chan bool {
	c := make(chan bool)

	go func() {
		err := client.cc.Close()
		if err != nil {
			client.logger.Errorf("shuttown tracking provider err: %w", err)
		}
		c <- true
	}()
	return c
}
