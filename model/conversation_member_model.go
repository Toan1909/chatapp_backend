package model

import "time"

type ConversationMember struct {
	UserId        string `json:"userId" db:"user_id, omitempty"`
	FullName      string `json:"fullName,omitempty" db:"fullname, omitempty"`
	Phone         string `json:"phone,omitempty" db:"phone, omitempty"`
	UrlProfilePic string `json:"urlProfilePic" db:"url_profile_pic, omitempty"`
	Status        bool   `json:"status" db:"status, omitempty"`
	JoinedAt        time.Time `json:"joinedAt,omitempty" db:"joined_at, omitempty"`
}