package entity

type MQTTACL struct {
	//gorm.Model
	Allow    bool    `gorm:"not null"`
	IPAddr   *string // 可为空
	Username *string // 可为空
	ClientID string  // 可为空
	Access   int     `gorm:"not null"` // 1=订阅, 2=发布, 3=订阅和发布
	Topic    string  `gorm:"not null"`
}

func (MQTTACL) TableName() string {
	return "mqtt_acl"
}
