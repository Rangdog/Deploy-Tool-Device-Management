package handler

import (
	"BE_Manage_device/config"
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/service"
	"BE_Manage_device/pkg"
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
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user   body    dto.UserRegisterRequest   true  "User Data"
// @Router       /api/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var user dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.UnknownError)
	}

	_, err := h.service.Register(user.FirstName, user.LastName, user.Password, user.Email)
	if err != nil {
		log.Error("Happened error when saving data to database. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusCreated, pkg.BuildReponse(constant.Success, ""))
}

// User godoc
// @Summary      Login
// @Description  Đăng nhập
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user   body    dto.UserLoginRequest   true  "User Data"
// @Router       /api/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var user dto.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}

	userLogged, accessToken, refreshToken, err := h.service.Login(user.Email, user.Password)
	if err != nil {
		log.Error("Happened error when login. Error", err)
		pkg.PanicExeption(constant.Invalidemailorpassword)
	}

	data := map[string]interface{}{
		"user": userLogged,
	}
	c.SetCookie("access_token", accessToken, 60*15, "/", config.BASE_URL_BACKEND, false, true)
	c.SetCookie("refresh_token", refreshToken, 60*60*24, "/", config.BASE_URL_BACKEND, false, true)

	c.JSON(http.StatusOK, pkg.BuildReponse(constant.Success, data))
}

func (h *UserHandler) Activate(c *gin.Context) {
	defer pkg.PanicHandler(c)
	token, exist := c.GetQuery("token")
	if !exist {
		log.Error("Happened error when mapping request from FE. Error: Dont see token in url")
		pkg.PanicExeption(constant.InvalidRequest)
	}
	err := h.service.Activate(token)
	if err != nil {
		log.Error("Happened error when activate. Error", err)
		pkg.PanicExeption(constant.UnknownError)
	}
	c.Redirect(http.StatusFound, config.BASE_URL_FRONTEND+"login")
}

// User godoc
// @Summary      Refresh Token
// @Description  Refresh Token
// @Tags         users
// @Accept       json
// @Produce      json
// @Router       /api/refresh [GET]
func (h *UserHandler) Refresh(c *gin.Context) {
	defer pkg.PanicHandler(c)
	refreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		log.Error("Happened error when refresh token. Error", err)
		pkg.PanicExeption(constant.Unauthorized)
	}
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
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
			pkg.PanicExeption(constant.Unauthorized)
			return
		}
		email := claims["email"].(string)
		user, err := h.service.FindUserByEmail(email)
		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": user.Id,
			"email":  email,
			"exp":    time.Now().Add(time.Minute * 15).Unix(),
		})

		newRefeshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": user.Id,
			"email":  email,
			"exp":    time.Now().Add(time.Hour * 14).Unix(),
		})

		if err != nil {
			pkg.PanicExeption(constant.UnknownError)
			return
		}

		newAccessTokenString, err := newAccessToken.SignedString([]byte(config.AccessSecret))
		if err != nil {
			pkg.PanicExeption(constant.UnknownError)
			return
		}
		newRefeshtokenString, err := newRefeshToken.SignedString([]byte(config.RefreshSecret))
		if err != nil {
			pkg.PanicExeption(constant.UnknownError)
			return
		}
		c.SetCookie("access_token", newAccessTokenString, 60*15, "/", config.BASE_URL_BACKEND, false, true)
		c.SetCookie("refresh_token", newRefeshtokenString, 60*60*24, "/", config.BASE_URL_BACKEND, false, true)

		c.JSON(http.StatusOK, pkg.BuildReponse(constant.Success, "Success"))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
	}
}
