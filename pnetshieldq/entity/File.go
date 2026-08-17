package entity

import "time"

type File struct {
	FileID     int
	FileName   string
	SenderID   string
	ReceiverID string
	FileURL    string
	Timestamp  time.Time
	FileType   string
}
