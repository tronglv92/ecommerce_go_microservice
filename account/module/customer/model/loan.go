package customermodel

import (
	"time"

	"github.com/tronglv92/accounts/common"
)

// CREATE TABLE `loans` (
//   `id` int NOT NULL AUTO_INCREMENT,
//   `customer_id` int NOT NULL,
//   `start_date` timestamp NOT NULL,
//   `loan_type` varchar(100) NOT NULL,
//   `total_loan` int NOT NULL,
//   `amount_paid` int NOT NULL,
//   `outstanding_amount` int NOT NULL,
//   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
//   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//   PRIMARY KEY (`id`)
// );

type Loan struct {
	common.SQLModel
	StartDate         *time.Time `json:"start_date" gorm:"column:start_date;"`
	LoanType          string     `json:"loan_type" gorm:"column:loan_type;"`
	TotalLoan         int        `json:"total_loan" gorm:"column:total_loan;"`
	AmountPaid        int        `json:"amount_paid" gorm:"column:amount_paid;"`
	OutStandingAmount int        `json:"outstanding_amount" gorm:"column:outstanding_amount;"`
}

func (r *Loan) Mask() {
	r.SQLModel.Mask(common.DbLoan)

}
