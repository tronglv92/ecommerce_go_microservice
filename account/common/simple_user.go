package common

type SimpleUser struct {
	SQLModel
	LastName   string `json:"last_name" gorm:"column:last_name;"`
	FirstName  string `json:"first_name" gorm:"column:first_name;"`
	Role       string `json:"role" gorm:"column:role;"`
	InternalId int    `json:"internal_id,omitempty"`

	//Avatar    *Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (SimpleUser) TableName() string {
	return "users"
}
