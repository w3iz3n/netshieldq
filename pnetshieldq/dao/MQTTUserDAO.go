package dao

import (
	"github.com/jinzhu/gorm"
	"pnetshieldq/entity"
)

type MQTTUserDAO struct {
	DB *gorm.DB
}

func NewMQTTUserDAO(db *gorm.DB) *MQTTUserDAO {
	return &MQTTUserDAO{DB: db}
}

func (dao *MQTTUserDAO) Create(user *entity.MqttUser) error {
	return dao.DB.Create(user).Error
}

func (dao *MQTTUserDAO) FindByUsername(username string) (*entity.MqttUser, error) {
	var user entity.MqttUser
	result := dao.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
