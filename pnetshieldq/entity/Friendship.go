package entity

import "time"

type Friendship struct {
	UserID1   int `gorm:"primaryKey"`
	UserID2   int `gorm:"primaryKey"`
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
