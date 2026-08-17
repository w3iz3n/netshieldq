package dao

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"pnetshieldq/entity"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}
func (dao *UserDao) Create(user *entity.User) error {
	user.SetLastSeen(time.Now())
	return dao.db.Create(user).Error
}

func (dao *UserDao) DeleteByUsername(username string) error {
	return dao.db.Where("username = ?", username).Delete(&entity.User{}).Error
}

func (dao *UserDao) Update(user *entity.User) error {
	// 使用 GORM 更新操作，只更新密码和邮箱
	err := dao.db.Model(user).Updates(map[string]interface{}{
		"email":     user.Email,
		"password":  user.Password, // 确保密码已经经过哈希处理
		"last_seen": user.LastSeen,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
func (dao *UserDao) UpdateFields(userID int, updates map[string]interface{}) error {
	return dao.db.Model(&entity.User{}).Where("user_id = ?", userID).Updates(updates).Error
}

func (dao *UserDao) Delete(user *entity.User) error {
	return dao.db.Delete(user).Error
}
func (dao *UserDao) DeleteByID(id int) error {
	return dao.db.Where("user_id = ?", id).Delete(&entity.User{}).Error
}
func (dao *UserDao) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	if err := dao.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (dao *UserDao) FindByID(id int) (*entity.User, error) {
	var user entity.User
	if err := dao.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := dao.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) UpdateClientID(userID int, clientID string) error {
	return dao.db.Model(&entity.User{}).Where("user_id = ?", userID).Update("client_id", clientID).Error
}

func (dao *UserDao) ResetPassword(username string, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return dao.db.Model(&entity.User{}).Where("username = ?", username).Update("password", string(hashedPassword)).Error
}

func (dao *UserDao) GetMaxUserID() (int, error) {
	var maxID int
	if err := dao.db.Model(&entity.User{}).Select("max(user_id)").Row().Scan(&maxID); err != nil {
		return 0, err
	}
	if maxID == 0 {
		return 1, nil
	}
	return maxID, nil
}
func (dao *UserDao) DeleteUser(id int) error {
	return dao.db.Where("id = ?", id).Delete(&entity.User{}).Error
}
