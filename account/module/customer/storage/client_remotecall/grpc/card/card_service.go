package remoterestful

import (
	"context"
	"log"

	"github.com/tronglv92/accounts/common"
	accountModel "github.com/tronglv92/accounts/module/customer/model"
	cardgrpc "github.com/tronglv92/accounts/proto/card"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

func (s *cardGrpcStore) GetCardsFromCustomerId(ctx context.Context, customerId int) ([]accountModel.Card, error) {
	logger := logger.GetCurrent().GetLogger("GetCardsFromCustomerId")
	res, err := s.client.ListCardByCustomerId(ctx, &cardgrpc.CardRequest{CustomerId: int32(customerId)})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	cards := []accountModel.Card{}

	for _, v := range res.Cards {
		logger.Debugf("card %v", v)
		createAt := v.GetCreatedAt().AsTime()
		updateAt := v.GetUpdateAt().AsTime()
		card := accountModel.Card{
			SQLModel: common.SQLModel{
				Id:        int(v.Id),
				CreatedAt: &createAt,
				UpdateAt:  &updateAt,
			},
			CardNumber:      v.CardNumber,
			CardType:        v.CardType,
			TotalLimit:      int(v.TotalLimit),
			AmountUsed:      int(v.AmountUsed),
			AvailableAmount: int(v.AvailableAmount),
		}
		card.Mask()
		cards = append(cards, card)

	}

	return cards, nil
}
