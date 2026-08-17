package dao

import (
	"github.com/jinzhu/gorm"
)

type AdminDao struct {
	db *gorm.DB
}

func NewAdminDao(db *gorm.DB) *AdminDao {
	return &AdminDao{db: db}
}

func (dao *AdminDao) CheckAdmin(username, password string) (bool, error) {
	var count int
	err := dao.db.Table("admins").
		Where("username = ? AND password = SHA2(?, 256)", username, password).
		Count(&count).Error
	return count > 0, err
}
