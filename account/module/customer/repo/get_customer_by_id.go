package customerrepo

import (
	"context"

	customermodel "github.com/tronglv92/accounts/module/customer/model"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

type GetCustomerByIdStorage interface {
	GetCustomer(
		context context.Context,
		cond map[string]interface{},
		moreKeys ...string) (*customermodel.Customer, error)
	GetAccountsFromCustomerId(
		context context.Context,
		customerId int,
		moreKeys ...string,
	) ([]customermodel.Account, error)
}

type GetCardsRemoteStorage interface {
	GetCardsFromCustomerId(
		context context.Context,
		id int,
	) ([]customermodel.Card, error)
}

type GetLoansRemoteStorage interface {
	GetLoansFromCustomerId(
		context context.Context,
		id int,
	) ([]customermodel.Loan, error)
}

type getCustomerRepo struct {
	store     GetCustomerByIdStorage
	cardStore GetCardsRemoteStorage
	loanStore GetLoansRemoteStorage
}

func NewCustomerByIdRepo(
	store GetCustomerByIdStorage,
	cardStore GetCardsRemoteStorage,
	loanStore GetLoansRemoteStorage,

) *getCustomerRepo {
	return &getCustomerRepo{
		store:     store,
		cardStore: cardStore,
		loanStore: loanStore,
	}
}
func (repo *getCustomerRepo) GetCustomerById(
	ctx context.Context,
	id int,
) (*customermodel.FullCustomer, error) {
	_ = logger.GetCurrent().GetLogger("account.repo.list_account")
	customer, err := repo.store.GetCustomer(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}

	cards, err := repo.cardStore.GetCardsFromCustomerId(ctx, id)
	if err != nil {
		return nil, err
	}

	accounts, err := repo.store.GetAccountsFromCustomerId(ctx, id)
	if err != nil {
		return nil, err
	}

	loans, err := repo.loanStore.GetLoansFromCustomerId(ctx, id)
	if err != nil {
		return nil, err
	}

	fullCustomer := customer.ConvertToFullCustomer(cards, accounts, loans)
	return fullCustomer, nil
}
