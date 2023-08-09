package remoterestful

import (
	cardgrpc "github.com/tronglv92/accounts/proto/card"
)

type cardGrpcStore struct {
	client cardgrpc.CardServiceClient
}

func NewCardRestfulStore(client cardgrpc.CardServiceClient) *cardGrpcStore {

	return &cardGrpcStore{client: client}
}
