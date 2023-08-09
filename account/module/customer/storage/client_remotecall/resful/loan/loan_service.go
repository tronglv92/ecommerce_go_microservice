package loanrestful

import (
	"context"
	"errors"
	"fmt"

	"log"

	accountModel "github.com/tronglv92/accounts/module/customer/model"
)

func (s *loanRestfulStore) GetLoansFromCustomerId(ctx context.Context, customerId int) ([]accountModel.Loan, error) {

	type responseCard struct {
		Data []accountModel.Loan `json:"data"`
	}

	var result responseCard

	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(fmt.Sprintf("%s/%s/%v", s.serviceURL, "internal/loans", customerId))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Println(resp.RawResponse)
		return nil, errors.New("cannot call api get loans")
	}

	for i := range result.Data {
		result.Data[i].GetRealId()
	}

	return result.Data, nil
}
