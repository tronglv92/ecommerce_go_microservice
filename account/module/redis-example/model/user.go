package usermodel

import "github.com/google/uuid"

const EntityName = "User"

type User struct {
	ID string `json:"id,omitempty"`

	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`

	// Avatar          *common.Image `json:"avatar" gorm:"column:avatar;type:json"`
}

func (data *User) Fullfill() {
	id := uuid.New()
	data.ID = id.String()
}
