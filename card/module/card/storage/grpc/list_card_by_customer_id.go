package grpcstorage

import (
	"context"
	"errors"

	cardmodel "github.com/tronglv92/cards/module/card/model"
	"github.com/tronglv92/ecommerce_go_common/logger"
	"google.golang.org/grpc/metadata"
)

func (s *sqlStore) ListCardByCustomerId(
	context context.Context,
	customerId int,
	moreKeys ...string) ([]cardmodel.Card, error) {
	logger := logger.GetCurrent().GetLogger("ListCardByCustomerIdStorage")
	// time.Sleep(5 * time.Second)
	var result []cardmodel.Card
	var empty []cardmodel.Card

	md, ok := metadata.FromIncomingContext(context)
	if !ok {
		return nil, errors.New("missing metadata")
	}
	logger.Debugf("md %v\n", md)
	s.gorm.WithContext(context)
	db := s.dbSession.Table(cardmodel.Card{}.TableName())

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
