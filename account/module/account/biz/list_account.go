package accountbiz

import (
	"context"

	"github.com/tronglv92/accounts/common"
	accountmodel "github.com/tronglv92/accounts/module/account/model"
)

type ListAccountRepo interface {
	ListAccount(
		context context.Context,
		filter *accountmodel.Filter,
		paging *common.Paging,
	) ([]accountmodel.Account, error)
}

type listAccountBiz struct {
	repo ListAccountRepo
}

func NewListAccountBiz(repo ListAccountRepo) *listAccountBiz {
	return &listAccountBiz{repo: repo}
}
func (biz *listAccountBiz) ListAccount(
	context context.Context,
	filter *accountmodel.Filter,
	paging *common.Paging,
) ([]accountmodel.Account, error) {
	result, err := biz.repo.ListAccount(context, filter, paging)
	if err != nil {
		return nil, common.ErrCannotListEntity(accountmodel.EntityName, err)
	}
	return result, nil
}
