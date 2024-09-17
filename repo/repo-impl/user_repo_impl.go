package repoimpl

import (
	"chatapp/db"
	my_err "chatapp/err"
	"chatapp/model"
	"chatapp/model/req"
	"chatapp/mylog"
	"chatapp/repo"
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type UserRepoImpl struct {
	sql *db.Sql
}

// GetUserInfo implements repo.UserRepo.
func (u *UserRepoImpl) GetUserInfo(c context.Context, userId string) (model.User, error) {
	user:= model.User{}
	statement:= `SELECT * FROM users where user_id = $1`
	err := u.sql.Db.GetContext(c,&user,statement,userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, my_err.UserNotFound
		}
		return user, err
	}
	return user,nil
}

// LoadListFriend implements repo.UserRepo.
func (u *UserRepoImpl) LoadListFriend(c context.Context, userId string) ([]model.User, error) {
	var listFriend []model.User
	statement := `SELECT u.user_id, u.fullname, u.phone, u.url_profile_pic,u.status
					FROM users u
					JOIN friendships f
					ON (u.user_id = f.friend_id AND f.user_id = $1)
					OR (u.user_id = f.user_id AND f.friend_id = $1)
					WHERE f.status = 'accepted';
				`
	err := u.sql.Db.SelectContext(c, &listFriend, statement, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return listFriend, my_err.FriendListNotFound
		}
		return listFriend, err
	}
	return listFriend, nil
}

// SaveFriShip implements repo.UserRepo.
func (u *UserRepoImpl) SaveFriShip(c context.Context, userId string, friendId string) (model.FriendShip, error) {
	fship := model.FriendShip{
		UserId:    userId,
		FriendId:  friendId,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	statement := `
		INSERT INTO 
			friendships(
				user_id,
				friend_id,
				status,
				created_at
			)
			VALUES(
				:user_id,
				:friend_id,
				:status,
				:created_at
			)
	`
	_, err := u.sql.Db.NamedExecContext(c, statement, fship)
	if err != nil {
		mylog.LogError(err)
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return fship, my_err.FriendshipConflict
			}
		}
	}
	return fship, nil
}

// SaveFriShip implements repo.UserRepo.

func NewUserRepo(sql *db.Sql) repo.UserRepo {
	return &UserRepoImpl{sql: sql}
}
func (u *UserRepoImpl) CheckLogIn(c context.Context, loginReq req.ReqSignIn) (model.User, error) {
	user := model.NewUser()
	err := u.sql.Db.GetContext(c, &user, "SELECT * FROM users WHERE email=$1", loginReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}
	return user, nil
}

// SaveUser implements repo.UserRepo.
func (u *UserRepoImpl) SaveUser(context context.Context, user model.User) (model.User, error) {
	statement := `
		INSERT INTO 
			users(
				user_id,
				fullname,
				email,
				phone,
				password,
				url_profile_pic,
				status,
				created_at
			)
			VALUES(
				:user_id,
				:fullname,
				:email,
				:phone,
				:password,
				:url_profile_pic,
				:status,
				:created_at
			)
	`
	user.CreatedAt = time.Now()
	_, err := u.sql.Db.NamedExecContext(context, statement, user)
	if err != nil {
		mylog.LogError(err)
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return user, my_err.UserConflict
			}
		}
	}
	return user, nil
}
