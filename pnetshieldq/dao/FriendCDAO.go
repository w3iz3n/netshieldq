package dao

import (
	"github.com/jinzhu/gorm"
	"pnetshieldq/entity"
)

type FriendCDAO struct {
	db *gorm.DB
}

func NewFriendCDAO(db *gorm.DB) *FriendCDAO {
	return &FriendCDAO{db: db}
}

func (dao *FriendCDAO) Create(friendC *entity.FriendC) error {
	return dao.db.Create(friendC).Error
}

func (dao *FriendCDAO) GetByID(id uint) (*entity.FriendC, error) {
	var friendC entity.FriendC
	result := dao.db.First(&friendC, id)
	return &friendC, result.Error
}
func (dao *FriendCDAO) GetByFriendIDAndUserID(userID uint, friendID uint) (*entity.FriendC, error) {
	var friendC entity.FriendC
	result := dao.db.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friendC)
	if result.Error != nil {
		return nil, result.Error
	}
	return &friendC, nil
}
func (dao *FriendCDAO) Update(friendC *entity.FriendC) error {
	return dao.db.Save(friendC).Error
}

func (dao *FriendCDAO) Delete(id uint) error {
	return dao.db.Delete(&entity.FriendC{}, id).Error
}
