package accountmodel

import (
	"time"

	"github.com/tronglv92/loans/common"
)

const EntityName = "loans"

type Loan struct {
	common.SQLModel
	StartDate         *time.Time `json:"start_date" gorm:"column:start_date;"`
	LoanType          string     `json:"loan_type" gorm:"column:loan_type;"`
	CustomerId        int        `json:"-" gorm:"column:customer_id;"`
	TotalLoan         int        `json:"total_loan" gorm:"column:total_loan;"`
	AmountPaid        int        `json:"amount_paid" gorm:"column:amount_paid;"`
	OutstandingAmount int        `json:"outstanding_amount" gorm:"column:outstanding_amount;"`

	FakeCustomerId *common.UID `json:"customer_id" gorm:"-"`
}

func (Loan) TableName() string { return "loans" }

func (r *Loan) Mask() {
	r.SQLModel.Mask(common.DbLoan)

	fakeCustomerId := common.NewUID(uint32(r.CustomerId), int(common.DbCustomer), 1)
	r.FakeCustomerId = &fakeCustomerId

}
