package accountmodel

type Filter struct {
	CustomerId int    `json:"customer_id,omitempty" form:"customer_id"`
	Search     string `json:"search,omitempty" form:"search"`
}
