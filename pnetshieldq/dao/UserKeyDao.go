package dao

import (
	"github.com/jinzhu/gorm"
	"pnetshieldq/entity"
)

type UserKeyDAO struct {
	db *gorm.DB
}

func NewUserKeyDAO(db *gorm.DB) *UserKeyDAO {
	return &UserKeyDAO{db: db}
}
func (dao *UserKeyDAO) UpsertKey(userKey *entity.UserKey) error {
	var existingUserKey entity.UserKey
	result := dao.db.Where("user_id = ?", userKey.UserID).First(&existingUserKey)

	if result.RecordNotFound() {
		return dao.db.Create(userKey).Error
	} else if result.Error != nil {
		return result.Error
	} else {
		updates := map[string]interface{}{
			"pk":         userKey.PK,
			"created_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}
		return dao.db.Model(&entity.UserKey{}).Where("user_id = ?", userKey.UserID).Updates(updates).Error
	}
}

func (dao *UserKeyDAO) CreateKey(userKey *entity.UserKey) error {
	return dao.db.Create(userKey).Error
}

func (dao *UserKeyDAO) RetrieveKeyByID(userID uint) (*entity.UserKey, error) {
	var userKey entity.UserKey
	result := dao.db.Where("user_id = ?", userID).First(&userKey)
	return &userKey, result.Error
}
func (dao *UserKeyDAO) ModifyKey(userKey *entity.UserKey) error {
	return dao.db.Save(userKey).Error
}

func (dao *UserKeyDAO) RemoveKey(userID uint) error {
	return dao.db.Delete(&entity.UserKey{}, userID).Error
}
func (dao *UserKeyDAO) GetAllKeys() ([]*entity.UserKey, error) {
	var userKeys []*entity.UserKey
	result := dao.db.Find(&userKeys)
	return userKeys, result.Error
}
