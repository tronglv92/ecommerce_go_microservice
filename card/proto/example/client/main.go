package main

import (
	"context"
	"log"
	"time"

	card "github.com/tronglv92/cards/proto/card"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.Dial("localhost:50051", opts)

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
