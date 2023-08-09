package accountbiz

import (
	"context"

	"github.com/tronglv92/cards/common"
	cardmodel "github.com/tronglv92/cards/module/card/model"
)

type ListCardRepo interface {
	ListCard(
		context context.Context,
		filter *cardmodel.Filter,
		paging *common.Paging,
	) ([]cardmodel.Card, error)
}

type listCardBiz struct {
	repo ListCardRepo
}

func NewListCardBiz(repo ListCardRepo) *listCardBiz {
	return &listCardBiz{repo: repo}
}
func (biz *listCardBiz) ListCard(
	context context.Context,
	filter *cardmodel.Filter,
	paging *common.Paging,
) ([]cardmodel.Card, error) {
	result, err := biz.repo.ListCard(context, filter, paging)
	if err != nil {
		return nil, common.ErrCannotListEntity(cardmodel.EntityName, err)
	}
	return result, nil
}
