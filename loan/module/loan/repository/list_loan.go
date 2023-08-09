package accountrepo

import (
	"context"

	"github.com/tronglv92/ecommerce_go_common/logger"
	"github.com/tronglv92/loans/common"
	loanmodel "github.com/tronglv92/loans/module/loan/model"
)

type ListLoanStorage interface {
	ListDataWithCondition(
		context context.Context,
		filter *loanmodel.Filter,
		paging *common.Paging,
		moreKeys ...string) ([]loanmodel.Loan, error)
}

type listLoanRepo struct {
	store ListLoanStorage

	// uStore UserStore
}

func NewListLoanRepo(store ListLoanStorage) *listLoanRepo {
	return &listLoanRepo{store: store}
}
func (repo *listLoanRepo) ListLoan(
	ctx context.Context,
	filter *loanmodel.Filter,
	paging *common.Paging,
) ([]loanmodel.Loan, error) {
	logger := logger.GetCurrent().GetLogger("loan.repo.list_loan")
	result, err := repo.store.ListDataWithCondition(ctx, filter, paging)
	if err != nil {
		return nil, err
	}
	logger.Debugf("vao trong nay", result)

	return result, nil
}
