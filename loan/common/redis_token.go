package common

import "time"

type RedisToken struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`

	//Avatar    *Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}
