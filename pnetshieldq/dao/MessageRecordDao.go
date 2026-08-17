package dao

import (
	"github.com/jinzhu/gorm"
	"pnetshieldq/entity"
	_ "pnetshieldq/entity"
)

type MessageRecordDao struct {
	db *gorm.DB
}

func NewMessageRecordDao(db *gorm.DB) *MessageRecordDao {
	return &MessageRecordDao{db: db}
}

func (dao *MessageRecordDao) FindAll() ([]*entity.MessageRecord, error) {
	var records []*entity.MessageRecord
	rows, err := dao.db.Table("messages").Select("sender_id, receiver_id,content, timestamp").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record entity.MessageRecord
		var senderID, receiverID int
		err = rows.Scan(&senderID, &receiverID, &record.Content, &record.Timestamp)
		if err != nil {
			return nil, err
		}

		var sender entity.User
		err = dao.db.Table("users").Where("user_id = ?", senderID).First(&sender).Error
		if err != nil {
			return nil, err
		}
		record.Username = sender.Username

		var receiver entity.User
		err = dao.db.Table("users").Where("user_id = ?", receiverID).First(&receiver).Error
		if err != nil {
			return nil, err
		}
		record.FriendName = receiver.Username

		records = append(records, &record)
	}

	return records, nil
}
