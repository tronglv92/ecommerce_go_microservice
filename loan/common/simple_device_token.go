package common

type SimpleDeviceToken struct {
	SQLModel
	Token        string `json:"token" gorm:"column:token;"`
	UserId       int    `json:"user_id" gorm:"column:user_id;"`
	DeviceId     string `json:"device_id" gorm:"column:device_id;"`
	IsProduction int    `json:"is_production" gorm:"column:is_production;default:0"`
	OS           string `json:"os" gorm:"column:os;"`
}

func (SimpleDeviceToken) TableName() string {
	return "user_device_tokens"
}
