package accountstorage

import (
	"context"

	loanmodel "github.com/tronglv92/loans/module/loan/model"
)

func (s *sqlStore) ListLoanByCustomerId(
	context context.Context,
	customerId int,
	moreKeys ...string) ([]loanmodel.Loan, error) {

	var result []loanmodel.Loan
	var empty []loanmodel.Loan

	db := s.db.Table(loanmodel.Loan{}.TableName())

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
