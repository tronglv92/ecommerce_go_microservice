package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	card "github.com/tronglv92/cards/proto/card"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct{}

func (s *server) ListCardByCustomerId(ctx context.Context, request *card.CardRequest) (*card.CardResponse, error) {
	log.Println(request.CustomerId)

	return &card.CardResponse{
		Cards: []*card.Card{
			{
				Id:              1,
				CustomerId:      1,
				CardNumber:      "1",
				CardType:        "1",
				TotalLimit:      1,
				AmountUsed:      1,
				AvailableAmount: 1,
			},
		},
	}, nil
}
func main() {
	address := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Error %v", err)
	}

	fmt.Printf("Server is listening on %v ...", address)

	s := grpc.NewServer()

	card.RegisterCardServiceServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:50051",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = card.RegisterCardServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":3000",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:3000")
	log.Fatalln(gwServer.ListenAndServe())
}
