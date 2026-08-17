package dao

import (
	"github.com/jinzhu/gorm"
	"pnetshieldq/entity"
)

type MQTTACLDAO struct {
	DB *gorm.DB
}

func NewMQTTACLDAO(db *gorm.DB) *MQTTACLDAO {
	return &MQTTACLDAO{DB: db}
}

func (dao *MQTTACLDAO) Insert(acl *entity.MQTTACL) error {
	return dao.DB.Create(acl).Error
}

func (dao *MQTTACLDAO) Delete(id uint) error {
	return dao.DB.Delete(&entity.MQTTACL{}, id).Error
}

func (dao *MQTTACLDAO) Update(acl *entity.MQTTACL) error {
	return dao.DB.Save(acl).Error
}

// FindByID 根据 ID 查找 ACL 记录
func (dao *MQTTACLDAO) FindByID(id uint) (*entity.MQTTACL, error) {
	var acl entity.MQTTACL
	result := dao.DB.First(&acl, id)
	return &acl, result.Error
}

// FindAll 获取所有 ACL 记录
func (dao *MQTTACLDAO) FindAll() ([]entity.MQTTACL, error) {
	var acls []entity.MQTTACL
	result := dao.DB.Find(&acls)
	return acls, result.Error
}
