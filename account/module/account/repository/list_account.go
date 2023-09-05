package accountrepo

import (
	"context"

	"github.com/tronglv92/accounts/common"
	accountmodel "github.com/tronglv92/accounts/module/account/model"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

type ListAccountStorage interface {
	ListDataWithCondition(
		context context.Context,
		filter *accountmodel.Filter,
		paging *common.Paging,
		moreKeys ...string) ([]accountmodel.Account, error)
}
type CustomerRepo interface {
	GetCustomersByIds(ctx context.Context, ids []int) ([]common.Customer, error)
}

type listAccountRepo struct {
	store         ListAccountStorage
	customerStore CustomerRepo
	// uStore UserStore
}

func NewListAccountRepo(store ListAccountStorage, customerStore CustomerRepo) *listAccountRepo {
	return &listAccountRepo{store: store, customerStore: customerStore}
}
func (repo *listAccountRepo) ListAccount(
	ctx context.Context,
	filter *accountmodel.Filter,
	paging *common.Paging,
) ([]accountmodel.Account, error) {
	_ = logger.GetCurrent().GetLogger("account.repo.list_account")
	result, err := repo.store.ListDataWithCondition(ctx, filter, paging)
	if err != nil {
		return nil, err
	}

	// fmt.Println("restaurant.repo.list_restaurant: ", result)
	customerIds := make([]int, len(result))

	for i := range customerIds {
		customerIds[i] = result[i].CustomerId
	}

	customers, err := repo.customerStore.GetCustomersByIds(ctx, customerIds)

	if err != nil {
		return nil, common.ErrCannotListEntity(accountmodel.EntityName, err)
	}
	mapUser := make(map[int]*common.Customer)

	for j, u := range customers {
		mapUser[u.Id] = &customers[j]
	}

	for i, item := range result {
		result[i].Customer = mapUser[item.CustomerId]
	}

	return result, nil
}
