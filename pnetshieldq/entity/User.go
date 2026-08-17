package entity

import "time"

const layout = "2006-01-02 15:04:05"

type User struct {
	UserID       int
	Username     string
	Password     string
	OnlineStatus bool
	LastSeen     time.Time
	IP           string
	Email        string
	ClientID     string
}

func (u *User) GetUserID() int {
	return u.UserID
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) GetOnlineStatus() bool {
	return u.OnlineStatus
}

func (u *User) GetLastSeen() string {
	return u.LastSeen.Format(layout)
}

func (u *User) GetIP() string {
	return u.IP
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetClientID() string {
	return u.ClientID
}

// Setters
func (u *User) SetUserID(userID int) {
	u.UserID = userID
}

func (u *User) SetUsername(username string) {
	u.Username = username
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

func (u *User) SetOnlineStatus(onlineStatus bool) {
	u.OnlineStatus = onlineStatus
}

func (u *User) SetLastSeen(lastSeen time.Time) {
	u.LastSeen = lastSeen
}

func (u *User) SetIP(ip string) {
	u.IP = ip
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) SetClientID(clientID string) {
	u.ClientID = clientID
}
