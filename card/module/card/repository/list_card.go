package accountrepo

import (
	"context"

	"github.com/tronglv92/cards/common"
	cardmodel "github.com/tronglv92/cards/module/card/model"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

type ListCardStorage interface {
	ListDataWithCondition(
		context context.Context,
		filter *cardmodel.Filter,
		paging *common.Paging,
		moreKeys ...string) ([]cardmodel.Card, error)
}

type listCardRepo struct {
	store ListCardStorage

	// uStore UserStore
}

func NewListCardRepo(store ListCardStorage) *listCardRepo {
	return &listCardRepo{store: store}
}
func (repo *listCardRepo) ListCard(
	ctx context.Context,
	filter *cardmodel.Filter,
	paging *common.Paging,
) ([]cardmodel.Card, error) {
	logger := logger.GetCurrent().GetLogger("card.repo.list_card")
	result, err := repo.store.ListDataWithCondition(ctx, filter, paging)
	if err != nil {
		return nil, err
	}
	logger.Debugf("vao trong nay", result)

	

	return result, nil
}
