package accountrepo

import (
	"context"

	"github.com/tronglv92/ecommerce_go_common/logger"
	loanmodel "github.com/tronglv92/loans/module/loan/model"
)

type ListLoanByCustomerIdStorage interface {
	ListLoanByCustomerId(
		context context.Context,
		customerId int,
		moreKeys ...string) ([]loanmodel.Loan, error)
}

type listLoanByCustomerId struct {
	store ListLoanByCustomerIdStorage

	// uStore UserStore
}

func NewListLoanByCustomerIdRepo(store ListLoanByCustomerIdStorage) *listLoanByCustomerId {
	return &listLoanByCustomerId{store: store}
}
func (repo *listLoanByCustomerId) ListLoanByCustomerId(
	ctx context.Context,
	customerId int,

) ([]loanmodel.Loan, error) {
	logger := logger.GetCurrent().GetLogger("loan.repo.list_loan")
	result, err := repo.store.ListLoanByCustomerId(ctx, customerId)
	if err != nil {
		return nil, err
	}
	logger.Debugf("vao trong nay", result)

	return result, nil
}
