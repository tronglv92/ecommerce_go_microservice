package accountbiz

import (
	"context"

	"github.com/tronglv92/loans/common"
	loanmodel "github.com/tronglv92/loans/module/loan/model"
)

type ListLoanByCustomerIdRepo interface {
	ListLoanByCustomerId(
		context context.Context,
		customerId int,
	) ([]loanmodel.Loan, error)
}

type listLoanByCustomerIdBiz struct {
	repo ListLoanByCustomerIdRepo
}

func NewListLoanByCustomerIdBiz(repo ListLoanByCustomerIdRepo) *listLoanByCustomerIdBiz {
	return &listLoanByCustomerIdBiz{repo: repo}
}
func (biz *listLoanByCustomerIdBiz) ListLoanByCustomerId(
	context context.Context,
	customerId int,
) ([]loanmodel.Loan, error) {

	result, err := biz.repo.ListLoanByCustomerId(context, customerId)
	if err != nil {
		return nil, common.ErrCannotListEntity(loanmodel.EntityName, err)
	}
	return result, nil
}
