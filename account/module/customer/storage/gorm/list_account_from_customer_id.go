package customerstore

import (
	"context"

	customermodel "github.com/tronglv92/accounts/module/customer/model"
)

func (s *sqlStore) GetAccountsFromCustomerId(
	context context.Context,
	customerId int,
	moreKeys ...string,
) ([]customermodel.Account, error) {

	var result []customermodel.Account
	var empty []customermodel.Account
	db := s.db.Table(customermodel.Account{}.TableName())

	if customerId > 0 {
		db = db.Where("customer_id=?", customerId)
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}
	if err := db.
		Order("id desc").
		Find(&result).Error; err != nil {
		return empty, err
	}

	return result, nil
}
