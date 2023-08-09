package accountstorage

import (
	"context"
	"strings"

	"github.com/tronglv92/accounts/common"
	accountmodel "github.com/tronglv92/accounts/module/account/model"
)

func (s *sqlStore) ListDataWithCondition(
	context context.Context,
	filter *accountmodel.Filter,
	paging *common.Paging,
	moreKeys ...string) ([]accountmodel.Account, error) {

	var result []accountmodel.Account
	var empty []accountmodel.Account

	db := s.db.Table(accountmodel.Account{}.TableName())

	if f := filter; f != nil {
		if f.CustomerId > 0 {
			db = db.Where("customer_id=?", f.CustomerId)
		}
		search := strings.TrimSpace(f.Search)
		if len(search) > 0 {
			db = db.Where("account_number=?", search)
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
		last.Mask(false)
		paging.NextCursor = last.FakeId.String()
	}

	return result, nil
}
