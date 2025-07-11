package utils

import (
	"BE_Manage_device/config"
	"BE_Manage_device/internal/domain/entity"
	"strconv"

	"BE_Manage_device/pkg/interfaces"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"

	asset_log "BE_Manage_device/internal/repository/asset_log"
	repository "BE_Manage_device/internal/repository/asset_log"
	asset "BE_Manage_device/internal/repository/assets"
	user "BE_Manage_device/internal/repository/user"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

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

func GenQrAndUpdate(repo asset.AssetsRepository, assetId int64, url string) {
	qrUrl, err := GenerateAssetQR(assetId, url)
	if err != nil {
		logrus.Info("Error when create qrurl")
		return
	}
	err = repo.UpdateQrURL(assetId, qrUrl)
	if err != nil {
		logrus.Info("Error when update qrurl into asset")
	}
}

type notificationJob struct {
	Emails  []string
	Subject string
	Body    string
}

func CheckAndSenMaintenanceNotification(db *gorm.DB, emailNotifier interfaces.EmailNotifier, assetRepo asset.AssetsRepository, userRepo user.UserRepository, notification interfaces.Notification, assetLogRepo repository.AssetsLogRepository) {
	loc, _ := time.LoadLocation("Asia/Bangkok")

	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)
	var schedules []entity.MaintenanceSchedules
	err := db.Where("start_date <= ? and start_date >= ?", endOfDay, startOfDay).Preload("Asset").Preload("Asset.OnwerUser").Find(&schedules).Error
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
			var userManagerAsset *entity.Users
			userManagerAsset, _ = userRepo.GetUserAssetManageOfDepartment(s.Asset.DepartmentId)
			userRoleAdmin, _ := userRepo.GetUserRoleAdmin()
			users := []*entity.Users{s.Asset.OnwerUser, userManagerAsset}
			users = append(users, userRoleAdmin...)
			if len(users) == 0 {
				log.Printf("⚠️ No users with notification permission for asset ID %d", s.AssetId)
				return nil // hoặc có thể return error nếu muốn rollback transaction
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
				NotifyDate: time.Now(),
			}
			if err := tx.Create(&notify).Error; err != nil {
				return fmt.Errorf("error inserting notification: %w", err)
			}
			assetLog := entity.AssetLog{
				AssetId:       asset.Id,
				ChangeSummary: fmt.Sprintf("Asset %d has started maintenance", asset.Id),
				Timestamp:     time.Now(),
				Action:        "Maintenance",
				CompanyId:     asset.CompanyId,
			}
			if _, err := assetLogRepo.Create(&assetLog, tx); err != nil {
				return fmt.Errorf("error create asset log: %w", err)
			}
			job.Emails = emails
			job.Subject = subject
			job.Body = body
			message := fmt.Sprintf("The asset (ID: %v) moved to 'Under Maintenance'", s.AssetId)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("SendNotificationToUsers panic:", r)
					}
				}()
				notification.SendNotificationToUsers(users, message, s.Asset)
			}()
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

func UpdateStatusWhenFinishMaintenance(db *gorm.DB, assetRepo asset.AssetsRepository, userRepo user.UserRepository, notification interfaces.Notification, assetLogRepo asset_log.AssetsLogRepository) {
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
				userManagerAsset, _ := userRepo.GetUserAssetManageOfDepartment(a.DepartmentId)
				usersToNotifications := []*entity.Users{}
				if a.OnwerUser != nil {
					usersToNotifications = append(usersToNotifications, a.OnwerUser)
				}
				if userManagerAsset != nil {
					usersToNotifications = append(usersToNotifications, userManagerAsset)
				}
				message := fmt.Sprintf("The asset (ID: %v) moved to 'In Use'", a.Id)
				go func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Println("SendNotificationToUsers panic:", r)
						}
					}()
					notification.SendNotificationToUsers(usersToNotifications, message, *a)
				}()
			}
			assetLog := entity.AssetLog{
				AssetId:       a.Id,
				ChangeSummary: fmt.Sprintf("Asset %d has started maintenance", a.Id),
				Timestamp:     time.Now(),
				Action:        "Maintenance",
				CompanyId:     a.CompanyId,
			}
			if _, err := assetLogRepo.Create(&assetLog, assetLogRepo.GetDB()); err != nil {
				log.Printf("❌ Error create asset log ")
			}
		}
	}
}

func SendEmailsForWarrantyExpiry(db *gorm.DB, emailNotifier interfaces.EmailNotifier, notification interfaces.Notification, assetRepo asset.AssetsRepository, userRepo user.UserRepository) {
	assets, err := assetRepo.GetAssetsWasWarrantyExpiry()
	if err != nil {
		log.Printf("❌ Error fetching assets : %v", err)
		return
	}
	var jobs []notificationJob
	for _, a := range assets {
		var job notificationJob
		var userHeadDepart *entity.Users
		var userManagerAsset *entity.Users
		userManagerAsset, _ = userRepo.GetUserAssetManageOfDepartment(a.DepartmentId)
		userRoleAdmin, _ := userRepo.GetUserRoleAdmin()
		users := []*entity.Users{a.OnwerUser, userHeadDepart, userManagerAsset}
		users = append(users, userRoleAdmin...)
		if len(users) == 0 {
			log.Printf("⚠️ No users with notification permission for asset ID %d", a.Id)
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
		status := "Pending"

		notify := entity.Notifications{
			NotifyDate: &now,
			Type:       &typ,
			AssetId:    &assetId,
			Status:     &status,
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

		message := fmt.Sprintf("The asset (ID: %v) has just been Expired", a.Id)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("SendNotificationToUsers panic:", r)
				}
			}()
			notification.SendNotificationToUsers(users, message, *a)
		}()
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

func PtrInt64(i int64) *int64 {
	return &i
}

func GetSuffixEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func ParseInt64Ptr(s *string) (int64, error) {
	if s == nil {
		return 0, fmt.Errorf("nil string")
	}
	return strconv.ParseInt(*s, 10, 64)
}

func ParseStrToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return i, err
}

func ParseStrToBool(s string) (bool, error) {
	i, err := strconv.ParseBool(s)
	return i, err
}
