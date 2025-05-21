package utils

import (
	"BE_Manage_device/config"
	"BE_Manage_device/constant"
	"BE_Manage_device/pkg"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
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

type SupabaseUploader struct {
	ProjectRef string
	ApiKey     string
	Bucket     string
}

func NewSupabaseUploader() *SupabaseUploader {
	return &SupabaseUploader{
		ProjectRef: config.SUPABASE_PROJECT_REF, // hoặc gán thẳng chuỗi
		ApiKey:     config.SupabaseKey,          // hoặc gán thẳng
		Bucket:     "images",                    // Tên bucket
	}
}

func (s *SupabaseUploader) Upload(objectPath string, file multipart.File, contentType string) (string, error) {
	defer file.Close()

	// Tạo buffer để đọc toàn bộ file
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return "", err
	}

	// Gửi request PUT (Upload file)
	url := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/%s/%s", s.ProjectRef, s.Bucket, objectPath)

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.ApiKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprint(buf.Len()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", string(body))
	}

	// Trả về URL public nếu bucket public
	publicURL := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/public/%s/%s", s.ProjectRef, s.Bucket, objectPath)
	return publicURL, nil
}

func CleanTimezoneLabel(input string) string {
	// Tìm vị trí của " ("
	idx := strings.Index(input, " (")
	if idx != -1 {
		// Cắt từ đầu đến trước " ("
		return input[:idx]
	}
	// Nếu không có gì để cắt, trả về nguyên gốc
	return input
}
