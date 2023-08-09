package accountbiz

import (
	"context"

	"github.com/tronglv92/loans/common"
	loanmodel "github.com/tronglv92/loans/module/loan/model"
)

type ListLoanRepo interface {
	ListLoan(
		context context.Context,
		filter *loanmodel.Filter,
		paging *common.Paging,
	) ([]loanmodel.Loan, error)
}

type listLoanBiz struct {
	repo ListLoanRepo
}

func NewListLoanBiz(repo ListLoanRepo) *listLoanBiz {
	return &listLoanBiz{repo: repo}
}
func (biz *listLoanBiz) ListLoan(
	context context.Context,
	filter *loanmodel.Filter,
	paging *common.Paging,
) ([]loanmodel.Loan, error) {
	result, err := biz.repo.ListLoan(context, filter, paging)
	if err != nil {
		return nil, common.ErrCannotListEntity(loanmodel.EntityName, err)
	}
	return result, nil
}
