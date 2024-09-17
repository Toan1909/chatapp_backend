package handler

import (
	my_err "chatapp/err"
	"chatapp/model"
	"chatapp/model/req"
	"chatapp/mylog"
	"chatapp/repo"
	"chatapp/security"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserRepo repo.UserRepo
}
func (u *UserHandler) HandleSignUp(c echo.Context) error {
	req := req.ReqSignUp{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	//hash password
	hash := security.HashAndSalt([]byte(req.Password))
	userId, err := uuid.NewUUID()
	if err != nil {
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	user := model.User{
		UserId:   userId.String(),
		FullName: req.FullName,
		Email:    req.Email,
		Phone: req.Phone,
		Password: hash,
		UrlProfilePic: "",
		Status: false,
		Token:    "",
	}
	user, err = u.UserRepo.SaveUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusConflict, model.Response{
			StatusCode: http.StatusConflict,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	//gene token
	token, err := security.GenToken(user)
	if err != nil {
		mylog.LogError(err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	user.Token = token
	user.Password = "" //trước khi return user về ,modify passwd="" để omitempty ẩn passwd
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Sign-up user thành công",
		Data:       user,
	})

	
}

func (u *UserHandler) HandleSignIn(c echo.Context) error {
	req := req.ReqSignIn{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	user, err := u.UserRepo.CheckLogIn(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	//check pass
	isTruePwd := security.ComparePasswords(user.Password, []byte(req.Password))
	if !isTruePwd {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Sai pass ,đăng nhập thất bại!",
			Data:       nil,
		})
	}
	//gene token
	token, err := security.GenToken(user)
	if err != nil {
		mylog.LogError(err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	user.Token = token
	user.Password = "" //ẩn mật khẩu đi trước khi trả về
	//Không có lỗi => return user
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Đăng nhập thành công",
		Data:       user,
	})
}

func (u *UserHandler) HandleFriendship(c echo.Context) error {
	req := req.ReqFrShip{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	frship, err := u.UserRepo.SaveFriShip(c.Request().Context(),req.UserId,req.FriendId)
	if err != nil {
		if err == my_err.FriendshipConflict{
			return c.JSON(http.StatusConflict, model.Response{
				StatusCode: http.StatusConflict,
				Message:    err.Error(),
				Data:       nil,
			})
		}
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Friendship thành công",
		Data:       frship,
	})
}
func (u *UserHandler) HandleGetListFriend(c echo.Context) error {
	req := req.ReqLoadFriendList{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	frList, err := u.UserRepo.LoadListFriend(c.Request().Context(),req.UserId)
	if err != nil {
		if err == my_err.FriendListNotFound{
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Load listfriend thành công",
		Data:       frList,
	})
}