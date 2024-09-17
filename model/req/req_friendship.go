package req

type ReqFrShip struct {
	UserId    string    `json:"userId,omitempty" db:"user_id, omitempty" validate:"required"`
	FriendId  string    `json:"friendId,omitempty" db:"friend_id, omitempty" validate:"required"`
}
type ReqLoadFriendList struct{
	UserId string `json:"userId,omitempty" db:"user_id, omitempty" validate:"required"`
}