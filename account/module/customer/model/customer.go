package customermodel

import "github.com/tronglv92/accounts/common"

const EntityName = "customers"

type Customer struct {
	common.SQLModel
	Name         string `json:"name" gorm:"column:name;"`
	Email        string `json:"email" gorm:"column:email;"`
	MobileNumber string `json:"mobile_number" gorm:"column:mobile_number;"`
}

func (Customer) TableName() string { return "customers" }

func (r *Customer) Mask() {
	r.SQLModel.Mask(common.DbCustomer)

}
func (r *Customer) ConvertToFullCustomer(cards []Card, accounts []Account, loans []Loan) *FullCustomer {

	return &FullCustomer{
		Cards:        cards,
		Loans:        loans,
		Accounts:     accounts,
		Name:         r.Name,
		Email:        r.Email,
		MobileNumber: r.MobileNumber,
		SQLModel:     r.SQLModel,
	}
}

type FullCustomer struct {
	common.SQLModel
	Name         string    `json:"name" gorm:"column:name;"`
	Email        string    `json:"email" gorm:"column:email;"`
	MobileNumber string    `json:"mobile_number" gorm:"column:mobile_number;"`
	Accounts     []Account `json:"accounts" gorm:"-"`
	Loans        []Loan    `json:"loans" gorm:"-"`
	Cards        []Card    `json:"cards" gorm:"-"`
}

func (r *FullCustomer) Mask() {
	r.SQLModel.Mask(common.DbCustomer)

}
