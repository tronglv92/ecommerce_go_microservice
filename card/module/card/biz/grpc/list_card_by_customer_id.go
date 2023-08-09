package grpcstore

import (
	"context"

	cardmodel "github.com/tronglv92/cards/module/card/model"
	cardgrpc "github.com/tronglv92/cards/proto/card"
	"github.com/tronglv92/ecommerce_go_common/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListCardByCustomerIdStorage interface {
	// user/storage/gorm/get: GetUsers
	ListCardByCustomerId(
		context context.Context,
		customerId int,
	) ([]cardmodel.Card, error)
}
type listCardByCustomerIdBiz struct {
	dbStore ListCardByCustomerIdStorage
}

func NewListCardByCustomerIdBiz(dbStore ListCardByCustomerIdStorage) *listCardByCustomerIdBiz {
	return &listCardByCustomerIdBiz{dbStore: dbStore}
}

func (s *listCardByCustomerIdBiz) ListCardByCustomerId(ctx context.Context, request *cardgrpc.CardRequest) (*cardgrpc.CardResponse, error) {
	logger := logger.GetCurrent().GetLogger("ListCardByCustomerId")
	rs, err := s.dbStore.ListCardByCustomerId(ctx, int(request.GetCustomerId()))
	if err != nil {
		return nil, err
	}

	cards := make([]*cardgrpc.Card, len(rs))

	for i, item := range rs {

		logger.Debugf("item %v", item.Id)
		item.Mask()

		cards[i] = &cardgrpc.Card{
			Id:              int32(item.Id),
			CustomerId:      int32(item.CustomerId),
			CardNumber:      item.CardNumber,
			CardType:        item.CardType,
			TotalLimit:      int32(item.TotalLimit),
			AmountUsed:      int32(item.AmountUsed),
			AvailableAmount: int32(item.AvailableAmount),
			CreatedAt:       timestamppb.New(*item.CreatedAt),
			UpdateAt:        timestamppb.New(*item.UpdateAt),
		}
	}
	return &cardgrpc.CardResponse{Cards: cards}, nil
}
