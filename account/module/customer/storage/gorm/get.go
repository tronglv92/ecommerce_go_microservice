package customerstore

import (
	"context"

	customermodel "github.com/tronglv92/accounts/module/customer/model"

	"github.com/tronglv92/accounts/common"
	"gorm.io/gorm"
)

func (s *sqlStore) GetCustomer(
	context context.Context,
	condition map[string]interface{},
	moreKeys ...string) (*customermodel.Customer, error) {
		
	db := s.db.Table(customermodel.Customer{}.TableName())
	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}
	var data customermodel.Customer
	if err := db.Where(condition).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &data, nil
}
