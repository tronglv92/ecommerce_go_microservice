package accountmodel

import (
	"github.com/tronglv92/accounts/common"
)

const EntityName = "accounts"

type Account struct {
	common.SQLModel
	CustomerId     int              `json:"-" gorm:"column:customer_id;"`
	AccountNumber  int              `json:"account_number" gorm:"column:account_number;"`
	AccountType    string           `json:"account_type" gorm:"column:account_type;"`
	BranchAddress  string           `json:"branch_address" gorm:"column:branch_address;"`
	FakeCustomerId *common.UID      `json:"customer_id" gorm:"-"`
	Customer       *common.Customer `json:"customer" gorm:"preload:false;"`
}

func (Account) TableName() string { return "accounts" }

func (r *Account) Mask(isAdminOrOwner bool) {
	r.SQLModel.Mask(common.DbAccount)

	fakeCustomerId := common.NewUID(uint32(r.CustomerId), int(common.DbCustomer), 1)
	r.FakeCustomerId = &fakeCustomerId

	if v := r.Customer; v != nil {
		v.Mask(common.DbCustomer)
	}

}
