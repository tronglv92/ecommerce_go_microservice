package customerbiz

import (
	"context"

	"github.com/tronglv92/accounts/common"
	customermodel "github.com/tronglv92/accounts/module/customer/model"
)

type GetCustomerRepo interface {
	GetCustomerById(
		ctx context.Context,
		id int,
	) (*customermodel.FullCustomer, error)
}
type getCustomerBiz struct {
	repo GetCustomerRepo
}

func NewGetCustomerBiz(repo GetCustomerRepo) *getCustomerBiz {
	return &getCustomerBiz{repo: repo}
}
func (biz *getCustomerBiz) GetCustomerById(
	context context.Context,
	id int,
) (*customermodel.FullCustomer, error) {
	result, err := biz.repo.GetCustomerById(context, id)
	if err != nil {
		return nil, common.ErrCannotListEntity(customermodel.EntityName, err)
	}
	return result, nil
}
