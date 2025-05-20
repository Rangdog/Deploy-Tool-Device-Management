package utils

import (
	"BE_Manage_device/config"
	"BE_Manage_device/constant"
	"BE_Manage_device/pkg"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

func GetUserIdFromContext(c *gin.Context) int64 {
	userID, exists := c.Get("userID")
	if exists {
		log.Info("userID:", userID)
	} else {
		log.Error("Happened error when get userId from gin Context")
		pkg.PanicExeption(constant.UnknownError)
	}
	str := fmt.Sprint(userID)

	userIdConvert, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Error("Happened error when get userId in token. Error", err)
		pkg.PanicExeption(constant.UnknownError)
	}
	return userIdConvert
}

func LogEmailError(action string, to string, err error) {
	log.WithFields(log.Fields{
		"action": action,
		"to":     to,
		"error":  err,
	}).Error("❌ Gửi email thất bại")
}

func LogEmailSuccess(action string, to string) {
	log.WithFields(log.Fields{
		"action": action,
		"to":     to,
	}).Info("✅ Gửi email thành công")
}

func GenerateTokens(userId int64, email string) (string, string, error) {
	// Access Token (15 phút)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"exp":    time.Now().Add(10 * time.Hour).Unix(),
	})
	accessString, err := accessToken.SignedString([]byte(config.AccessSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token (7 ngày)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"exp":    time.Now().Add(50 * time.Hour).Unix(),
	})
	refreshString, err := refreshToken.SignedString([]byte(config.RefreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}
