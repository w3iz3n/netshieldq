package entity

import (
	"time"
)

type UserKey struct {
	UserID    uint      `gorm:"primaryKey;autoIncrement"`
	PK        string    `gorm:"type:text;notNull"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	status    string    `gorm:"default:'active"`
}
type FriendC struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"notNull"`
	FriendID  uint      `gorm:"notNull"`
	C         string    `gorm:"notNull"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
}
