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
	"github.com/golang-jwt/jwt/v4"
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
func (u *UserHandler) HandleGetProfile(c echo.Context) error{
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomclaims)
	user,err :=u.UserRepo.GetUserInfo(c.Request().Context(),claims.UserId)
	if err!= nil {
		if err==my_err.UserNotFound{
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    "Không tìm thấy thông tin người dùng",
				Data:       user,
			})
		}
		
		return c.JSON(http.StatusInternalServerError, model.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
			})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Fetch profile thành công",
		Data:       user,
	})
}
func (u *UserHandler) HandleGetProfileOther(c echo.Context) error{
	req := req.ReqProfile{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	user,err :=u.UserRepo.GetUserInfo(c.Request().Context(),req.UserId)
	if err!= nil {
		if err==my_err.UserNotFound{
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    "Không tìm thấy thông tin người dùng",
				Data:       user,
			})
		}
		
		return c.JSON(http.StatusInternalServerError, model.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
			})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Fetch profile thành công",
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
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomclaims)
	frship, err := u.UserRepo.SaveFriendShip(c.Request().Context(),claims.UserId,req.FriendId)
	
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
		Message:    "Addfriend thành công",
		Data:       frship,
	})
}
func (u *UserHandler) HandleAcceptFriendship(c echo.Context) error {
	req := req.ReqFrShip{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomclaims)
	err := u.UserRepo.AcceptFriendShip(c.Request().Context(),claims.UserId,req.FriendId)
	
	if err != nil {
		if err == my_err.FriendshipConflict{
			return c.JSON(http.StatusConflict, model.Response{
				StatusCode: http.StatusConflict,
				Message:    err.Error(),
				Data:       false,
			})
		}
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Accept thành công",
		Data:       true,
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
func (u *UserHandler) HandleGetListPending(c echo.Context) error {
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

	frList, err := u.UserRepo.LoadListPending(c.Request().Context(),req.UserId)
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
		Message:    "Load pendings thành công",
		Data:       frList,
	})
}
func (u *UserHandler) HandleCheckFriend(c echo.Context) error {
	req := req.ReqCheckFrShip{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomclaims)
	isFriend, _ := u.UserRepo.CheckFriend(c.Request().Context(),claims.UserId,req.FriendId)
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Check Friend thành công",
		Data:       isFriend,
	})
}


func (u *UserHandler) HandleSearchUser(c echo.Context) error {
	req := req.ReqSearchUser{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	

	user, err := u.UserRepo.SearchUser(c.Request().Context(),req.Email,req.Phone)
	if err != nil {
		if err == my_err.SearchNotFound{
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Search thành công",
		Data:       user,
	})
}