package customerstore

import (
	"context"

	"github.com/tronglv92/accounts/common"
)

func (s *sqlStore) GetCustomersByIds(ctx context.Context, ids []int) ([]common.Customer, error) {
	var result []common.Customer

	if err := s.db.Table(common.Customer{}.TableName()).
		Where("id in (?)", ids).
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}
