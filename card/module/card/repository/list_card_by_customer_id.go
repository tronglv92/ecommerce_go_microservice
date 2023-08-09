package accountrepo

import (
	"context"

	cardmodel "github.com/tronglv92/cards/module/card/model"
)

type ListCardByCustomerIdStorage interface {
	ListCardByCustomerId(
		context context.Context,
		customerId int,
		moreKeys ...string) ([]cardmodel.Card, error)
}

type listCardByCustomerIdRepo struct {
	store ListCardByCustomerIdStorage

	// uStore UserStore
}

func NewListCardByCustomerIdRepo(store ListCardByCustomerIdStorage) *listCardByCustomerIdRepo {
	return &listCardByCustomerIdRepo{store: store}
}
func (repo *listCardByCustomerIdRepo) ListCardByCustomerId(
	ctx context.Context,
	customerId int,

) ([]cardmodel.Card, error) {
	// logger := logger.GetCurrent().GetLogger("card.repo.list_card")
	result, err := repo.store.ListCardByCustomerId(ctx, customerId)
	if err != nil {
		return nil, err
	}

	return result, nil
}
