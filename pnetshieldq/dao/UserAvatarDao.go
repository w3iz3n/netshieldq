package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"pnetshieldq/entity"
)

type UserAvatarDAO struct {
	DB *gorm.DB
}

func NewUserAvatarDAO(db *gorm.DB) *UserAvatarDAO {
	return &UserAvatarDAO{DB: db}
}

func (dao *UserAvatarDAO) InsertUserAvatar(avatar string) (*entity.UserAvatar, error) {
	userAvatar := entity.UserAvatar{Avatar: avatar}
	if err := dao.DB.Create(&userAvatar).Error; err != nil {
		return nil, err
	}
	fmt.Printf("Inserted UserAvatar ID: %d\n", userAvatar.ID) // 打印插入后的 ID
	return &userAvatar, nil
}

func (dao *UserAvatarDAO) UpdateUserAvatar(userID int, avatar string) error {
	userAvatar := entity.UserAvatar{
		ID:     userID,
		Avatar: avatar,
	}

	// 使用 Save 方法更新头像
	if err := dao.DB.Model(&userAvatar).Where("id = ?", userID).Update("avatar", avatar).Error; err != nil {
		return err
	}

	return nil
}

func (dao *UserAvatarDAO) GetUserAvatar(userID int) ([]byte, error) {
	var userAvatar entity.UserAvatar

	if err := dao.DB.Where("id = ?", userID).First(&userAvatar).Error; err != nil {
		return nil, errors.New("user avatar not found")
	}

	return []byte(userAvatar.Avatar), nil
}
