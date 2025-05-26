package utils

import (
	"BE_Manage_device/config"
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/interfaces"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
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

func (s *SupabaseUploader) UploadReader(objectPath string, reader io.Reader, contentType string) (string, error) {
	// Đọc toàn bộ nội dung từ reader vào buffer
	var buf bytes.Buffer
	_, err := io.Copy(&buf, reader)
	if err != nil {
		return "", fmt.Errorf("failed to read content: %w", err)
	}

	// Tạo URL Supabase Storage
	url := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/%s/%s", s.ProjectRef, s.Bucket, objectPath)

	// Tạo request POST (upload)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.ApiKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprint(buf.Len()))

	// Gửi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	// Kiểm tra trạng thái phản hồi
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", string(body))
	}

	// Trả về URL public (nếu bucket là public)
	publicURL := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/public/%s/%s", s.ProjectRef, s.Bucket, objectPath)
	return publicURL, nil
}

func GenerateAssetQR(assetID int64) (string, error) {
	url := fmt.Sprintf("%s/assets/%d", config.BASE_URL_FRONTEND, assetID)
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("QR encoding failed: %w", err)
	}
	// Tạo tên file và đường dẫn
	fileName := fmt.Sprintf("qr_%d_%d.png", assetID, time.Now().UnixNano())
	path := "qr-codes/" + fileName
	// Tạo reader để upload
	reader := bytes.NewReader(png)
	contentType := "image/png"
	uploader := NewSupabaseUploader()
	qrURL, err := uploader.UploadReader(path, reader, contentType)
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	return qrURL, nil
}

func CheckAndSenMaintenanceNotification(db *gorm.DB, emailNotifier interfaces.EmailNotifier, assetRepo repository.AssetsRepository) {
	today := time.Now().Truncate(24 * time.Hour)
	var schedules []entity.MaintenanceSchedules
	err := db.Where("start_date <= ? AND end_date >= ?", today, today).Find(&schedules).Error
	if err != nil {
		log.Printf("Error fetching maintenance schedules: %v", err)
		return
	}
	for _, s := range schedules {
		// Check nếu đã thông báo rồi
		var noti entity.MaintenanceNotifications
		err := db.Where("schedule_id = ?", s.Id).First(&noti).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("✅ Already notified for schedule ID %d today", s.Id)
			continue
		}

		// Nếu chưa có thông báo thì tiến hành
		err = db.Transaction(func(tx *gorm.DB) error {
			// 1. Lấy user nhận email
			users, err := assetRepo.GetUserHavePermissionNotifications(s.AssetId)
			if len(users) == 0 {
				log.Printf("⚠️ No users with notification permission for asset ID %d", s.AssetId)
				return nil // hoặc có thể return error nếu muốn rollback transaction
			}
			if err != nil {
				return fmt.Errorf("error fetching emails: %w", err)
			}

			// 2. Lấy asset
			asset, err := assetRepo.GetAssetById(s.AssetId)
			if err != nil {
				return fmt.Errorf("error fetching asset: %w", err)
			}

			// 3. Chuẩn bị email
			var emails []string
			for _, u := range users {
				emails = append(emails, u.Email)
			}
			subject := fmt.Sprintf("Asset %s is scheduled for maintenance on %s", asset.AssetName, s.StartDate.Format("Jan 2, 2006"))
			body := fmt.Sprintf(`
			<html>
				<body>
					<p>Dear team,</p>
					<p>Please be informed that the following asset is scheduled for maintenance:</p>
					<table border="1" cellpadding="6" cellspacing="0" style="border-collapse: collapse;">
						<tr>
							<th align="left">Asset</th>
							<td>%s</td>
						</tr>
						<tr>
							<th align="left">Start Date</th>
							<td>%s</td>
						</tr>
						<tr>
							<th align="left">End Date</th>
							<td>%s</td>
						</tr>
					</table>
					<p>Kindly plan accordingly.</p>
					<p>Best regards,<br>Your Maintenance Team</p>
				</body>
			</html>
		`, asset.AssetName, s.StartDate.Format("Jan 2, 2006"), s.EndDate.Format("Jan 2, 2006"))

			// 5. Cập nhật lifecycle
			if _, err := assetRepo.UpdateAssetLifeCycleStage(asset.Id, "Under Maintenance", tx); err != nil {
				return fmt.Errorf("error updating asset stage: %w", err)
			}

			// 6. Ghi log notification
			notify := entity.MaintenanceNotifications{
				ScheduleId: s.Id,
				NotifyDate: today,
			}
			if err := tx.Create(&notify).Error; err != nil {
				return fmt.Errorf("error inserting notification: %w", err)
			}

			// 4. Gửi email (ngoài transaction, async)
			go emailNotifier.SendEmails(emails, subject, body)
			return nil
		})

		if err != nil {
			log.Printf("❌ Transaction failed for schedule %d: %v", s.Id, err)
		}
	}
}

func UpdateStatusWhenFinishMaintenance(db *gorm.DB, assetRepo repository.AssetsRepository) {
	assets, err := assetRepo.GetAssetByStatus("Under Maintenance")
	if err != nil {
		log.Printf("❌ Error fetching assets with status 'Under Maintenance': %v", err)
		return
	}

	for _, a := range assets {
		finished, err := assetRepo.CheckAssetFinishMaintenance(a.Id)
		if err != nil {
			log.Printf("⚠️ Error checking maintenance status for asset %d: %v", a.Id, err)
			continue
		}

		if finished {
			_, err := assetRepo.UpdateAssetLifeCycleStage(a.Id, "In Use", db)
			if err != nil {
				log.Printf("❌ Error updating asset %d to 'In Use': %v", a.Id, err)
			} else {
				log.Printf("✅ Asset %d moved to 'In Use'", a.Id)
			}
		}
	}
}
