package utils

import (
	"BE_Manage_device/config"
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/interfaces"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
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

func (s *SupabaseUploader) Delete(objectPath string) error {
	// Chuẩn bị body (Supabase yêu cầu mảng path)
	paths := []string{objectPath}
	body, err := json.Marshal(map[string]interface{}{
		"prefixes": paths,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/%s", s.ProjectRef, paths[0])

	req, err := http.NewRequest("DELETE", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call supabase: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(b))
	}

	return nil
}

func ExtractFilePath(url string) (string, bool) {
	sep := "/public/"
	idx := strings.Index(url, sep)
	if idx != -1 {
		return url[idx+len(sep):], true
	}
	return "", false
}

func GenerateAssetQR(assetID int64, urlFrontend string) (string, error) {
	url := fmt.Sprintf("%s/%d", urlFrontend, assetID)
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

type notificationJob struct {
	Emails  []string
	Subject string
	Body    string
}

func CheckAndSenMaintenanceNotification(db *gorm.DB, emailNotifier interfaces.EmailNotifier, assetRepo repository.AssetsRepository) {
	today := time.Now().Truncate(24 * time.Hour)
	var schedules []entity.MaintenanceSchedules
	err := db.Where("start_date <= ? AND end_date >= ?", today, today).Find(&schedules).Error
	if err != nil {
		log.Printf("Error fetching maintenance schedules: %v", err)
		return
	}
	var jobs []notificationJob
	for _, s := range schedules {
		// Check nếu đã thông báo rồi
		var noti entity.MaintenanceNotifications
		err := db.Where("schedule_id = ?", s.Id).First(&noti).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("✅ Already notified for schedule ID %d today", s.Id)
			continue
		}
		var job notificationJob
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
			job.Emails = emails
			job.Subject = subject
			job.Body = body
			return nil
		})
		if len(job.Emails) > 0 {
			jobs = append(jobs, job)
		}
		if err != nil {
			log.Printf("❌ Transaction failed for schedule %d: %v", s.Id, err)
		}
	}
	const workerCount = 10
	jobsQueue := make(chan notificationJob, len(jobs))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsQueue {
				emailNotifier.SendEmails(job.Emails, job.Subject, job.Body)
			}
		}()
	}
	for _, job := range jobs {
		jobsQueue <- job
	}
	close(jobsQueue)
	wg.Wait()
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

func SendEmailsForWarrantyExpiry(db *gorm.DB, emailNotifier interfaces.EmailNotifier, assetRepo repository.AssetsRepository) {
	assets, err := assetRepo.GetAssetsWasWarrantyExpiry()
	if err != nil {
		log.Printf("❌ Error fetching assets : %v", err)
		return
	}
	var jobs []notificationJob
	for _, a := range assets {
		var job notificationJob
		users, err := assetRepo.GetUserHavePermissionNotifications(a.Id)
		if len(users) == 0 {
			log.Printf("⚠️ No users with notification permission for asset ID %d", a.Id)
			continue
		}
		if err != nil {
			log.Printf("❌ error fetching emails %d", a.Id)
			continue
		}
		var emails []string
		for _, u := range users {
			emails = append(emails, u.Email)
		}
		subject := fmt.Sprintf("Asset %s is expired on %s", a.AssetName, a.WarrantExpiry.Format("Jan 2, 2006"))
		body := fmt.Sprintf(`
			<html>
				<body>
					<p>Dear team,</p>
					<p>Please be informed that the following asset is expired:</p>
					<table border="1" cellpadding="6" cellspacing="0" style="border-collapse: collapse;">
						<tr>
							<th align="left">Asset</th>
							<td>%s</td>
						</tr>
						<tr>
							<th align="left">Expiry Date</th>
							<td>%s</td>
						</tr>
					</table>
					<p>Kindly plan accordingly.</p>
					<p>Best regards,<br>Your Manager Asset Team</p>
				</body>
			</html>
		`, a.AssetName, a.WarrantExpiry.Format("Jan 2, 2006"))
		now := time.Now()
		typ := "Expired"
		assetId := a.Id

		notify := entity.Notifications{
			NotifyDate: &now,
			Type:       &typ,
			AssetId:    &assetId,
		}
		result := db.Create(&notify)
		if result.Error != nil {
			log.Infof("Happen error when create notify type %v AssetId %v", typ, assetId)
		}
		job.Emails = emails
		job.Subject = subject
		job.Body = body
		if len(emails) > 0 {
			jobs = append(jobs, job)
		}
	}
	const workerCount = 10
	jobsQueue := make(chan notificationJob, len(jobs))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsQueue {
				emailNotifier.SendEmails(job.Emails, job.Subject, job.Body)
			}
		}()
	}
	for _, job := range jobs {
		jobsQueue <- job
	}
	close(jobsQueue)
	wg.Wait()
}

func UserHasPermission(db *gorm.DB, userId int64, permSlug []string, accessLevel []*string) (bool, error) {
	var user entity.Users
	err := db.Preload("Role.RolePermissions.Permission").
		First(&user, userId).Error
	if err != nil {
		return false, err
	}

	for _, rolePerm := range user.Role.RolePermissions {
		if accessLevel == nil {
			if slices.Contains(permSlug, rolePerm.Permission.Slug) && rolePerm.AccessLevel == "full" {
				return true, nil
			}
		} else {
			if slices.Contains(permSlug, rolePerm.Permission.Slug) && slices.Contains(accessLevel, &rolePerm.AccessLevel) {
				return true, nil
			}
		}
	}
	return false, nil
}

// Hàm tính giá trị còn lại của tài sản
func CurrentAssetValue(
	originalCost float64,
	salvageValue float64,
	usefulLifeYears float64,
	startDate, currentDate time.Time,
) float64 {

	// Tính số năm đã sử dụng (dạng float, có thể lẻ tháng)
	yearsUsed := currentDate.Sub(startDate).Hours() / (24 * 365)

	if yearsUsed < 0 {
		yearsUsed = 0 // chưa đến ngày bắt đầu sử dụng
	}
	if yearsUsed > usefulLifeYears {
		yearsUsed = usefulLifeYears // không vượt quá thời gian sử dụng
	}

	annualDepreciation := (originalCost - salvageValue) / usefulLifeYears
	accumulatedDepreciation := yearsUsed * annualDepreciation

	// Giá trị còn lại không thấp hơn salvage value
	currentValue := math.Max(originalCost-accumulatedDepreciation, salvageValue)
	return currentValue
}

func UpdateCost(assetRepo repository.AssetsRepository) {
	assets, err := assetRepo.GetAllAsset()
	if err != nil {
		log.Info("Happen error when get all asset")
	}
	for _, asset := range assets {
		if asset.Status != "New" || asset.Status != "Disposed" || asset.Status != "Retired" {
			if asset.AcquisitionDate != nil {
				cost := CurrentAssetValue(asset.OriginalCost, asset.ResidualValue, asset.UsefulLife, *asset.AcquisitionDate, time.Now())
				err = assetRepo.UpdateCost(asset.Id, cost)
				if err != nil {
					log.Info("Happen error when update cost")
				}
			}
		}
	}
}
