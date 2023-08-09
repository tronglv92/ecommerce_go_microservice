package customermodel

import (
	"github.com/tronglv92/accounts/common"
)

// CREATE TABLE `cards` (
//   `id` int NOT NULL AUTO_INCREMENT,
//   `card_number` varchar(100) NOT NULL,
//   `customer_id` int NOT NULL,
//   `card_type` varchar(100) NOT NULL,
//   `total_limit` int NOT NULL,
//   `amount_used` int NOT NULL,
//   `available_amount` int NOT NULL,
//    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
//   `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//   PRIMARY KEY (`id`)
// );

type Card struct {
	common.SQLModel

	CardNumber      string `json:"card_number" gorm:"column:card_number;"`
	CardType        string `json:"card_type" gorm:"column:card_type;"`
	TotalLimit      int    `json:"total_limit" gorm:"column:total_limit;"`
	AmountUsed      int    `json:"amount_used" gorm:"column:amount_used;"`
	AvailableAmount int    `json:"available_amount" gorm:"column:available_amount;"`
}

func (r *Card) Mask() {
	r.SQLModel.Mask(common.DbCard)

}
