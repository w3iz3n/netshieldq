package entity

type MessageRecord struct {
	Username   string `json:"username"`
	FriendName string `json:"friend_name"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
}
