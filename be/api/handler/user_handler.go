package handler

import (
	"BE_Manage_device/config"
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// User godoc
// @Summary      Register user
// @Description  Đăng ký user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user   body    dto.UserRegisterRequest   true  "Data"
// @Router       /api/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var user dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.UnknownError)
	}
	_, err := h.service.Register(user.FirstName, user.LastName, user.Password, user.Email, user.RedirectUrl)
	if err != nil {
		log.Error("Happened error when saving data to database. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusCreated, pkg.BuildReponse(constant.Success, ""))
}

// User godoc
// @Summary      Login
// @Description  Đăng nhập
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user   body    dto.UserLoginRequest   true  "Data"
// @Router       /api/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var user dto.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}

	_, accessToken, refreshToken, err := h.service.Login(user.Email, user.Password)
	if err != nil {
		log.Error("Happened error when login. Error", err)
		pkg.PanicExeption(constant.Invalidemailorpassword)
	}
	dataResponese := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	c.JSON(http.StatusOK, pkg.BuildReponse(constant.Success, dataResponese))
}

func (h *UserHandler) Activate(c *gin.Context) {
	defer pkg.PanicHandler(c)
	token, exist := c.GetQuery("token")
	if !exist {
		log.Error("Happened error when mapping request from FE. Error: Dont see token in url")
		pkg.PanicExeption(constant.InvalidRequest)
	}
	redirectUrl, exist := c.GetQuery("redirectUrl")
	if !exist {
		log.Error("Happened error when mapping request from FE. Error: Dont see token in url")
		pkg.PanicExeption(constant.InvalidRequest)
	}
	err := h.service.Activate(token)
	if err != nil {
		log.Error("Happened error when activate. Error", err)
		pkg.PanicExeption(constant.UnknownError)
	}
	c.Redirect(http.StatusFound, redirectUrl)
}

// User godoc
// @Summary      Refresh Token
// @Description  Refresh Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token   body    dto.RefreshRequest   true  "Data"
// @Router       /api/auth/refresh [POST]
func (h *UserHandler) Refresh(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var rq dto.RefreshRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}

	ok := h.service.CheckRefreshToken(rq.RefreshToken)
	if !ok {
		log.Error("Happened error refresh token was invoked")
		pkg.PanicExeption(constant.Unauthorized, "Refresh token was invoked")
	}
	refreshToken, err := jwt.Parse(rq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(config.RefreshSecret), nil
	})
	if err != nil || !refreshToken.Valid {
		pkg.PanicExeption(constant.Unauthorized)
		return
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		if !ok {
			pkg.PanicExeption(constant.Unauthorized)
			return
		}
		exp, ok := claims["exp"].(float64)
		if !ok {
			pkg.PanicExeption(constant.Unauthorized)
			return
		}
		if int64(exp) < time.Now().Unix() {
			pkg.PanicExeption(constant.StatusForbidden, "Refresh token was expired")
			return
		}
		email := claims["email"].(string)
		user, err := h.service.FindUserByEmail(email)
		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": user.Id,
			"email":  email,
			"exp":    time.Now().Add(time.Minute * 1).Unix(),
		})

		newRefeshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": user.Id,
			"email":  email,
			"exp":    time.Now().Add(time.Minute * 5).Unix(),
		})

		if err != nil {
			pkg.PanicExeption(constant.UnknownError)
		}

		newAccessTokenString, err := newAccessToken.SignedString([]byte(config.AccessSecret))
		if err != nil {
			pkg.PanicExeption(constant.UnknownError)
		}
		newRefeshtokenString, err := newRefeshToken.SignedString([]byte(config.RefreshSecret))
		if err != nil {
			pkg.PanicExeption(constant.UnknownError)
		}
		data := map[string]interface{}{
			"access_token":  newAccessTokenString,
			"refresh_token": newRefeshtokenString,
		}
		c.JSON(http.StatusOK, pkg.BuildReponse(constant.Success, data))
	} else {
		pkg.PanicExeption(constant.Unauthorized, "invalid refresh token")
	}
}

// User godoc
// @Summary      Reset Password
// @Description  Đặt lại mật khẩu
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        Reset_Password   body    dto.UserRequestResetPassword   true  "Data"
// @Router       /api/user/email/password-reset [PATCH]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.UserRequestResetPassword
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	user, err := h.service.FindUserByEmail(request.Email)
	if err != nil {
		log.Error("Happened error when email don't exist. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	err = h.service.ResetPassword(user, request.NewPassword, request.OldPassword)
	if err != nil {
		log.Error("Happened error when resert password. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponse(constant.Success, ""))
}

// User godoc
// @Summary      Get session
// @Description  Get session
// @Tags         users
// @Accept       json
// @Produce      json
// @Router       /api/user/session [GET]
func (h *UserHandler) Session(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	user, err := h.service.FindByUserId(userId)
	if err != nil {
		log.Error("Happened error when reset password. Error", err)
		pkg.PanicExeption(constant.Unauthorized, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponse(constant.Success, user))
}

// User godoc
// @Summary      Email reset password
// @Description   Email reset password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        Email_Reset_Password   body    dto.CheckPasswordReset   true  "Data"
// @Router       /api/user/password-reset [POST]
func (h *UserHandler) CheckPasswordReset(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.CheckPasswordReset
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	err := h.service.CheckPasswordReset(request.Email, request.RedirectUrl)
	if err != nil {
		log.Error("Happened error when reset password. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponse(constant.Success, ""))
}

// User godoc
// @Summary      Delete user
// @Description   Delete user via emal
// @Tags         users
// @Accept       json
// @Produce      json
// @Param		email	path		string				true	"email"
// @Router       /api/user/{email} [DELETE]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	email := c.Param("email")
	err := h.service.DeleteUser(email)
	if err != nil {
		log.Error("Happened error when delete user. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponse(constant.Success, ""))
}

// User godoc
// @Summary      Logout
// @Description   Logout
// @Tags         auth
// @Accept       json
// @Produce      json
// @Router       /api/auth/logout [POST]
func (h *UserHandler) Logout(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	userSession, err := h.service.FindSessionById(userId)
	if err != nil {
		log.Error("Happened error when logout user. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	err = h.service.UpdateInvoked(*userSession)
	if err != nil {
		log.Error("Happened error when logout user. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponse(constant.Success, ""))
}
