package remoterestful

import (
	"github.com/go-resty/resty/v2"
)

type cardRestfulStore struct {
	client     *resty.Client
	serviceURL string
}

func NewCardRestfulStore(client *resty.Client, serviceURL string) *cardRestfulStore {

	return &cardRestfulStore{client: client, serviceURL: serviceURL}
}
