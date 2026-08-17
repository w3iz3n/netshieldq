package dao

import (
	"github.com/jinzhu/gorm"
	"pnetshieldq/entity"
)

type FileDao struct {
	db *gorm.DB
}

func NewFileDao(db *gorm.DB) *FileDao {
	return &FileDao{db: db}
}

func (dao *FileDao) FindByName(name string) (*entity.File, error) {
	var file entity.File
	if err := dao.db.Where("file_name = ?", name).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (dao *FileDao) Create(file *entity.File) error {
	return dao.db.Create(file).Error
}

func (dao *FileDao) Update(file *entity.File) error {
	return dao.db.Save(file).Error
}

func (dao *FileDao) Delete(file *entity.File) error {
	return dao.db.Delete(file).Error
}

func (dao *FileDao) FindAll() ([]*entity.File, error) {
	var files []*entity.File
	if err := dao.db.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (dao *FileDao) FindBySenderID(senderID string) ([]*entity.File, error) {
	var files []*entity.File
	if err := dao.db.Where("sender_id = ?", senderID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
func (dao *FileDao) GetByUserID(userID string) ([]*entity.File, error) {
	var files []*entity.File
	if err := dao.db.Where("sender_id = ? OR receiver_id = ?", userID, userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
func (dao *FileDao) FindByReceiverID(receiverID string) ([]*entity.File, error) {
	var files []*entity.File
	if err := dao.db.Where("receiver_id = ?", receiverID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
func (dao *FileDao) FindByID(fileID string) (*entity.File, error) {
	var file entity.File
	if err := dao.db.Where("id = ?", fileID).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}
