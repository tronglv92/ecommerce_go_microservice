package accountstorage

import (
	"context"

	"github.com/tronglv92/cards/common"
	cardmodel "github.com/tronglv92/cards/module/card/model"
)

func (s *sqlStore) ListDataWithCondition(
	context context.Context,
	filter *cardmodel.Filter,
	paging *common.Paging,
	moreKeys ...string) ([]cardmodel.Card, error) {

	var result []cardmodel.Card
	var empty []cardmodel.Card

	db := s.db.Table(cardmodel.Card{}.TableName())

	if f := filter; f != nil {
		if f.CustomerId > 0 {
			db = db.Where("customer_id=?", f.CustomerId)
		}

	}

	if err := db.Count(&paging.Total).Error; err != nil {

		return empty, err
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)
		if err != nil {
			return empty, common.ErrDB(err)
		}
		db = db.Where("id<?", uid.GetLocalID())
	} else {
		offset := (paging.Page - 1) * paging.Limit
		db = db.Offset(offset)
	}

	if err := db.
		Limit(paging.Limit).
		Order("id desc").
		Find(&result).Error; err != nil {
		return empty, err
	}

	if len(result) > 0 {
		last := result[len(result)-1]
		last.Mask()
		paging.NextCursor = last.FakeId.String()
	}

	return result, nil
}
