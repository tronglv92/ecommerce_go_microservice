package main

import (
	"context"
	"log"

	"time"

	card "github.com/tronglv92/cards/proto/card"
	"github.com/tronglv92/cards/proto/example/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	tp, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.Dial("localhost:50051", opts,
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()
	client := card.NewCardServiceClient(cc)

	for i := 1; i <= 5; i++ {
		res, err := client.ListCardByCustomerId(context.Background(), &card.CardRequest{CustomerId: 1})

		if err != nil {
			log.Println(err)
		} else {
			log.Println(res.Cards)
		}

		time.Sleep(time.Second * 6)
	}
}
