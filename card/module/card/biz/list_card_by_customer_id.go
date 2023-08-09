package accountbiz

import (
	"context"

	"github.com/tronglv92/cards/common"
	cardmodel "github.com/tronglv92/cards/module/card/model"
)

type ListCardByCustomerIdRepo interface {
	ListCardByCustomerId(
		context context.Context,
		customerId int,
	) ([]cardmodel.Card, error)
}

type listCardByCustomerIdBiz struct {
	repo ListCardByCustomerIdRepo
}

func NewListCardByCustomerIdBiz(repo ListCardByCustomerIdRepo) *listCardByCustomerIdBiz {
	return &listCardByCustomerIdBiz{repo: repo}
}
func (biz *listCardByCustomerIdBiz) ListCardByCustomerId(
	context context.Context,
	customerId int,

) ([]cardmodel.Card, error) {
	result, err := biz.repo.ListCardByCustomerId(context, customerId)
	if err != nil {
		return nil, common.ErrCannotListEntity(cardmodel.EntityName, err)
	}
	return result, nil
}
