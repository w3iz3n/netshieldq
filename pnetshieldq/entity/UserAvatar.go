package entity

type UserAvatar struct {
	ID     int    `gorm:"primary_key;auto_increment"`
	Avatar string `gorm:"type:text;not null"`
}
