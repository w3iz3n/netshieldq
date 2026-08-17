package dao

import (
	"github.com/jinzhu/gorm"
	"pnetshieldq/entity"
)

type MessageDao struct {
	db *gorm.DB
}

func NewMessageDao(db *gorm.DB) *MessageDao {
	return &MessageDao{db: db}
}

func (dao *MessageDao) FindByID(id string) (*entity.Message, error) {
	var message entity.Message
	if err := dao.db.Where("message_id = ?", id).First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (dao *MessageDao) Create(message *entity.Message) error {
	return dao.db.Create(message).Error
}

func (dao *MessageDao) Update(message *entity.Message) error {
	return dao.db.Save(message).Error
}

func (dao *MessageDao) Delete(message *entity.Message) error {
	return dao.db.Delete(message).Error
}

func (dao *MessageDao) FindAll() ([]*entity.Message, error) {
	var messages []*entity.Message
	if err := dao.db.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

type MessageStats struct {
	Hour      int `json:"hour"`
	TextCount int `json:"text_count"`
	FileCount int `json:"file_count"`
}

func (dao *MessageDao) CountMessageTypesByHour() ([]MessageStats, error) {
	var results []MessageStats
	query := `SELECT
        HOUR(CONVERT_TZ(timestamp, @@session.time_zone, '+00:00')) as hour,
        SUM(CASE WHEN message_type = 'text' THEN 1 ELSE 0 END) as text_count,
        SUM(CASE WHEN message_type = 'file' THEN 1 ELSE 0 END) as file_count
    FROM
        messages
    WHERE
        DATE(CONVERT_TZ(timestamp, @@session.time_zone, '+00:00')) BETWEEN CURDATE() - INTERVAL 1 DAY AND CURDATE()
    GROUP BY
        hour
    ORDER BY
        hour;`

	if err := dao.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (dao *MessageDao) FindByMessageType(messageType string) ([]*entity.Message, error) {
	var messages []*entity.Message
	if err := dao.db.Where("message_type = ?", messageType).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (dao *MessageDao) FindByFileName(FileName string) ([]*entity.Message, error) {
	var messages []*entity.Message
	if err := dao.db.Where("file_name = ?", FileName).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
func (dao *MessageDao) GetMessagesBetweenUsers(userID uint, friendID uint) ([]*entity.Message, error) {
	var messages []*entity.Message
	if err := dao.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, friendID, friendID, userID).Order("timestamp desc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
