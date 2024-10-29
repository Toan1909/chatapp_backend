package req

type ReqFrShip struct {
	FriendId  string    `json:"friendId,omitempty" db:"friend_id, omitempty" validate:"required"`
}
type ReqCheckFrShip struct {
	FriendId  string    `json:"friendId,omitempty" db:"friend_id, omitempty" validate:"required"`
}
type ReqLoadFriendList struct{
	UserId string `json:"userId,omitempty" db:"user_id, omitempty" validate:"required"`
}