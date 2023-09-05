package customerbiz

import (
	"github.com/gin-gonic/gin"
	"github.com/tronglv92/accounts/common"
	customermodel "github.com/tronglv92/accounts/module/customer/model"
)

type GetCustomerRepo interface {
	GetCustomerById(
		c *gin.Context,
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
	c *gin.Context,
	id int,
) (*customermodel.FullCustomer, error) {
	result, err := biz.repo.GetCustomerById(c, id)
	if err != nil {
		return nil, common.ErrCannotListEntity(customermodel.EntityName, err)
	}
	return result, nil
}
