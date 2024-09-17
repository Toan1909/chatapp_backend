package my_err
import (
	"errors"
)
var (
	UserConflict = errors.New("Người dùng đã tồn tại")
	SignUpFail = errors.New("Đăng kí người dùng thất bại")
	UserNotFound =errors.New("Không tìm thấy người dùng/Không tồn tại")
	UserUpdateFail =errors.New("Cập nhật thông tin người dùng thất bại")
	MemberConflict = errors.New("Member đã tồn tại")
	ConvsersNotFound =errors.New("Không tìm thấy cuộc trò chuyện nào")
	MemNotFound =errors.New("Không tìm thấy cuộc trò chuyện nào")
	MessageNotFound =errors.New("Không tìm thấy tin nhắn nào")
	FriendshipConflict =errors.New("Friendship đã tồn tại")
	FriendListNotFound =errors.New("Không tìm thấy bạn bè nào")
)