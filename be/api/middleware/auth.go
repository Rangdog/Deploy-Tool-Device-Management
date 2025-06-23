package middleware

import (
	"BE_Manage_device/constant"
	repository "BE_Manage_device/internal/repository/user_session"

	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AuthMiddleware(secretKey string, session repository.UsersSessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer pkg.PanicHandler(c)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.PanicExeption(constant.Unauthorized, "Unauthorized Access Token")
			c.Abort()
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			pkg.PanicExeption(constant.Unauthorized, "Access Token expired")
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					pkg.PanicExeption(constant.Unauthorized, "Access Token expired")
					c.Abort()
					return
				}
			}
			a := claims["userId"]
			logrus.Info(a)
			c.Set("userID", claims["userId"])
			userID, exists := c.Get("userID")
			if !exists {
				pkg.PanicExeption(constant.UnknownError, "Happened error add auth")
			}
			str := fmt.Sprint(userID)

			userIdConvert, _ := strconv.ParseInt(str, 10, 64)

			if !session.CheckUserInSession(userIdConvert) {
				pkg.PanicExeption(constant.Unauthorized, "Unauthorized Access Token")
				c.Abort()
				return
			}
			if session.CheckTokenWasInVoked(tokenString) {
				pkg.PanicExeption(constant.Unauthorized, "Access Token was revoked")
				c.Abort()
				return
			}
		} else {
			pkg.PanicExeption(constant.Unauthorized, "Unauthorized Access Token")
			c.Abort()
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RequirePermission(permSlug []string, accessLevel []string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer pkg.PanicHandler(c)
		userID, exists := c.Get("userID")
		if !exists {
			pkg.PanicExeption(constant.Unauthorized, "Unauthorized Access Token")
			c.Abort()
			return
		}
		if exists {
			logrus.Info("userID:", userID)
		} else {
			logrus.Error("Happened error when get userId from gin Context")
			pkg.PanicExeption(constant.UnknownError, "Happened error when get userId from gin Context")
		}
		str := fmt.Sprint(userID)

		userIdConvert, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			pkg.PanicExeption(constant.UnknownError, "Internal server error")
			c.Abort()
			return
		}
		ok, err := utils.UserHasPermission(db, userIdConvert, permSlug, accessLevel)
		if err != nil {
			pkg.PanicExeption(constant.UnknownError, "Internal server error")
			c.Abort()
			return
		}
		if !ok {
			pkg.PanicExeption(constant.StatusForbidden, "Forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}

func TimeoutMiddleware(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
