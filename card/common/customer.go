package common

type Customer struct {
	SQLModel
	Name         string `json:"name" gorm:"column:name;"`
	Email        string `json:"email" gorm:"column:email;"`
	MobileNumber string `json:"mobile_number" gorm:"column:mobile_number;"`

	//Avatar    *Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (Customer) TableName() string {
	return "customers"
}
