package entity

import (
	"time"
)

type Message struct {
	MessageID   uint      `gorm:"primary_key;auto_increment" json:"message_id"`
	SenderID    uint      `gorm:"not null" json:"sender_id"`
	ReceiverID  uint      `gorm:"not null" json:"receiver_id"`
	Content     string    `gorm:"type:text" json:"content"`
	MessageType string    `gorm:"type:varchar(255)" json:"message_type"`
	Filename    string    `gorm:"type:text" json:"filename"`
	Timestamp   time.Time `gorm:"default:current_timestamp" json:"timestamp"`
	Status      string    `gorm:"type:varchar(20);default:'unread'" json:"status"`
}

func (m *Message) GetMessageType() string {
	return m.MessageType
}

func (m *Message) SetMessageType(messageType string) {
	m.MessageType = messageType
}

func (m *Message) GetFileName() string {
	return m.Filename
}

func (m *Message) setFileName(FileName string) {
	m.Filename = FileName
}
func (m *Message) GetTimestamp() string {
	return m.Timestamp.Format(layout)
}
