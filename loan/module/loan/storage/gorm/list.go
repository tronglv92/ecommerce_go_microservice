package accountstorage

import (
	"context"

	"github.com/tronglv92/loans/common"
	loanmodel "github.com/tronglv92/loans/module/loan/model"
)

func (s *sqlStore) ListDataWithCondition(
	context context.Context,
	filter *loanmodel.Filter,
	paging *common.Paging,
	moreKeys ...string) ([]loanmodel.Loan, error) {

	var result []loanmodel.Loan
	var empty []loanmodel.Loan

	db := s.db.Table(loanmodel.Loan{}.TableName())

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
