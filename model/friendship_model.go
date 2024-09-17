package model

import "time"

type FriendShip struct {
	UserId    string    `json:"userId,omitempty" db:"user_id, omitempty" validate:"required"`
	FriendId  string    `json:"friendId,omitempty" db:"friend_id, omitempty" validate:"required"`
	Status    string    `json:"status,omitempty" db:"status, omitempty" validate:"required"` //: pending, accepted, blocked
	CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at, omitempty" validate:"required"`
}