package customermodel

import (
	"github.com/tronglv92/accounts/common"
)

type Account struct {
	common.SQLModel
	AccountNumber int    `json:"account_number" gorm:"column:account_number;"`
	AccountType   string `json:"account_type" gorm:"column:account_type;"`
	BranchAddress string `json:"branch_address" gorm:"column:branch_address;"`
}

func (Account) TableName() string { return "accounts" }

func (r *Account) Mask() {
	r.SQLModel.Mask(common.DbAccount)

}
