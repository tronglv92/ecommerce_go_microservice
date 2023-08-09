package accountmodel

import (
	"github.com/tronglv92/cards/common"
)

const EntityName = "cards"

type Card struct {
	common.SQLModel
	CustomerId      int         `json:"-" gorm:"column:customer_id;"`
	CardNumber      string      `json:"card_number" gorm:"column:card_number;"`
	CardType        string      `json:"card_type" gorm:"column:card_type;"`
	TotalLimit      int         `json:"total_limit" gorm:"column:total_limit;"`
	AmountUsed      int         `json:"amount_used" gorm:"column:amount_used;"`
	AvailableAmount int         `json:"available_amount" gorm:"column:available_amount;"`
	FakeCustomerId  *common.UID `json:"customer_id" gorm:"-"`
}

func (Card) TableName() string { return "cards" }

func (r *Card) Mask() {
	r.SQLModel.Mask(common.DbCard)

	fakeCustomerId := common.NewUID(uint32(r.CustomerId), int(common.DbCustomer), 1)
	r.FakeCustomerId = &fakeCustomerId

}
