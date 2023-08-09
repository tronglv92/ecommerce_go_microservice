package accountmodel

type Filter struct {
	CustomerId int `json:"customer_id,omitempty" form:"customer_id"`
}
