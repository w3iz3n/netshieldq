package entity

type MqttUser struct {
	//gorm.Model
	Username    string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	IsSuperuser bool   `gorm:"default:false"`
}

func (MqttUser) TableName() string {
	return "mqtt_user"
}
