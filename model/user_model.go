package model

import "time"

type User struct {
	UserId   string    `json:"userId,omitempty" db:"user_id, omitempty"`
	FullName string    `json:"fullName,omitempty" db:"fullname, omitempty"`
	Email    string    `json:"email,omitempty" db:"email, omitempty"`
	Phone    string    `json:"phone,omitempty" db:"phone, omitempty"`
	Password string    `json:"password,omitempty" db:"password, omitempty"`
	UrlProfilePic string `json:"urlProfilePic" db:"url_profile_pic, omitempty"`
	Status bool `json:"status" db:"status, omitempty"`
	CreatedAt time.Time `json:"-" db:"created_at, omitempty"`
	Token    string    `json:"token,omitempty"`
}
func NewUser() User {
	return User{}
}
