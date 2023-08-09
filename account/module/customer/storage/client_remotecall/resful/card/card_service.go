package remoterestful

import (
	"context"
	"errors"
	"fmt"

	"log"

	accountModel "github.com/tronglv92/accounts/module/customer/model"
)

func (s *cardRestfulStore) GetCardsFromCustomerId(ctx context.Context, customerId int) ([]accountModel.Card, error) {

	type responseCard struct {
		Data []accountModel.Card `json:"data"`
	}

	var result responseCard

	url := fmt.Sprintf("%s/%s/%v", s.serviceURL, "internal/cards", customerId)

	fmt.Printf("url==%v", url)
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(url)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Println(resp.RawResponse)
		return nil, errors.New("cannot call api get cards")
	}

	for i := range result.Data {
		result.Data[i].GetRealId()
	}

	return result.Data, nil
}
