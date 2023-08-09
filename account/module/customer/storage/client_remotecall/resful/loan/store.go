package loanrestful

import (
	"github.com/go-resty/resty/v2"
)

type loanRestfulStore struct {
	client     *resty.Client
	serviceURL string
}

func NewLoanRestfulStore(client *resty.Client, serviceURL string) *loanRestfulStore {

	return &loanRestfulStore{client: client, serviceURL: serviceURL}
}
