package cronjob

import (
	asset_log "BE_Manage_device/internal/repository/asset_log"
	asset "BE_Manage_device/internal/repository/assets"
	bill "BE_Manage_device/internal/repository/bill"
	company "BE_Manage_device/internal/repository/company"
	monthlySummary "BE_Manage_device/internal/repository/monthly_summary"
	user "BE_Manage_device/internal/repository/user"
	emailS "BE_Manage_device/internal/service/email"
	notificationS "BE_Manage_device/internal/service/notification"
	"BE_Manage_device/pkg/utils"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func InitCronJobs(db *gorm.DB, emailService *emailS.EmailService, assetsRepository asset.AssetsRepository, userRepository user.UserRepository, notificationsService *notificationS.NotificationService, assetsLogRepository asset_log.AssetsLogRepository, billRepository bill.BillsRepository, monthlySummaryRepository monthlySummary.MonthlySummaryRepository, companyRepository company.CompanyRepository) {
	c := cron.New(cron.WithLocation(time.FixedZone("Asia/Ho_Chi_Minh", 7*3600)))

	_, err := c.AddFunc("0 8 * * *", func() {
		log.Println("🔔 Running maintenance notification check at 8:00 AM")
		utils.CheckAndSenMaintenanceNotification(db, emailService, assetsRepository, userRepository, notificationsService, assetsLogRepository)
	})

	if err != nil {
		log.Fatalf("❌ Failed to schedule maintenance cron job: %v", err)
	}

	_, err = c.AddFunc("1 8 * * *", func() {
		log.Println("🔔 Running warranty notification check at 8:00 AM")
		utils.SendEmailsForWarrantyExpiry(db, emailService, notificationsService, assetsRepository, userRepository)
	})
	if err != nil {
		log.Fatalf("❌ Failed to schedule warranty cron job: %v", err)
	}

	_, err = c.AddFunc("0 9 * * *", func() {
		log.Println("🔔 Running update status when finish maintenance at 9:00 AM")
		utils.UpdateStatusWhenFinishMaintenance(db, assetsRepository, userRepository, notificationsService, assetsLogRepository)
	})
	if err != nil {
		log.Fatalf("❌ Failed to schedule update status cron job: %v", err)
	}

	_, err = c.AddFunc("*/10 * * * *", func() {
		log.Println("🔔 Running kill session")
		utils.KillIdleSessions(db)
	})
	if err != nil {
		log.Fatalf("❌ Failed to schedule kill session cron job: %v", err)
	}

	_, err = c.AddFunc("0 0 * * *", func() {
		now := time.Now()
		// Get tomorrow's date
		tomorrow := now.AddDate(0, 0, 1)
		// If tomorrow is the first, today is the last day
		if tomorrow.Day() == 1 {
			utils.CreateSummary(monthlySummaryRepository, billRepository, companyRepository)
		} else {
			fmt.Println("Not the last day of the month. Skipping task...")
		}
	})
	if err != nil {
		log.Fatalf("❌ Failed to schedule create monthly summary cron job: %v", err)
	}

	c.Start()
}
