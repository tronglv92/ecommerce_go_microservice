package accountstorage

import (
	"context"

	cardmodel "github.com/tronglv92/cards/module/card/model"
)

func (s *sqlStore) ListCardByCustomerId(
	context context.Context,
	customerId int,
	moreKeys ...string) ([]cardmodel.Card, error) {

	var result []cardmodel.Card
	var empty []cardmodel.Card

	db := s.db.Table(cardmodel.Card{}.TableName())

	db = db.Where("customer_id=?", customerId)

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
