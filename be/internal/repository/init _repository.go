package repository

import (
	asset_log "BE_Manage_device/internal/repository/asset_log"
	asset "BE_Manage_device/internal/repository/assets"
	assignment "BE_Manage_device/internal/repository/assignments"
	bill "BE_Manage_device/internal/repository/bill"
	categories "BE_Manage_device/internal/repository/categories"
	company "BE_Manage_device/internal/repository/company"
	department "BE_Manage_device/internal/repository/departments"
	location "BE_Manage_device/internal/repository/locations"
	maintenanceNotification "BE_Manage_device/internal/repository/maintenance_notifications"
	maintenanceSchedules "BE_Manage_device/internal/repository/maintenance_schedules"
	monthlySummary "BE_Manage_device/internal/repository/monthly_summary"
	notification "BE_Manage_device/internal/repository/noftifications"
	request_transfer "BE_Manage_device/internal/repository/request_transfer"
	role "BE_Manage_device/internal/repository/role"
	user "BE_Manage_device/internal/repository/user"
	userRBAC "BE_Manage_device/internal/repository/user_rbac"
	userSession "BE_Manage_device/internal/repository/user_session"

	"gorm.io/gorm"
)

type Repository struct {
	User                    user.UserRepository
	UserSession             userSession.UsersSessionRepository
	Location                location.LocationRepository
	Categories              categories.CategoriesRepository
	Department              department.DepartmentsRepository
	Assets                  asset.AssetsRepository
	AssetsLog               asset_log.AssetsLogRepository
	Role                    role.RoleRepository
	UserRBAC                userRBAC.UserRBACRepository
	Assignment              assignment.AssignmentRepository
	RequestTransfer         request_transfer.RequestTransferRepository
	MaintenanceSchedules    maintenanceSchedules.MaintenanceSchedulesRepository
	Notification            notification.NotificationRepository
	MaintenanceNotification maintenanceNotification.MaintenanceNotificationsRepository
	Company                 company.CompanyRepository
	Bill                    bill.BillsRepository
	MonthlySummary          monthlySummary.MonthlySummaryRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:                    user.NewPostgreSQLUserRepository(db),
		UserSession:             userSession.NewPostgreSQLUserSessionRepository(db),
		Location:                location.NewPostgreSQLLocationRepository(db),
		Categories:              categories.NewPostgreSQLCategoriesRepository(db),
		Department:              department.NewPostgreSQLDepartmentsRepository(db),
		Assets:                  asset.NewPostgreSQLAssetsRepository(db),
		AssetsLog:               asset_log.NewPostgreSQLAssetsLogRepository(db),
		Role:                    role.NewPostgreSQLRoleRepository(db),
		UserRBAC:                userRBAC.NewPostgreSQLUserRBACRepository(db),
		Assignment:              assignment.NewPostgreSQLAssignmentRepository(db),
		RequestTransfer:         request_transfer.NewPostgreSQLRequestTransferRepository(db),
		MaintenanceSchedules:    maintenanceSchedules.NewPostgreSQLMaintenanceSchedulesRepository(db),
		Notification:            notification.NewPostgreSQLNotificationRepository(db),
		MaintenanceNotification: maintenanceNotification.NewPostgreSQLMaintenanceNotificationRepository(db),
		Company:                 company.NewPostgreSQLCompanyRepository(db),
		Bill:                    bill.NewPostgreSQLBillsRepository(db),
		MonthlySummary:          monthlySummary.NewPostgreSQLMonthlySummary(db),
	}
}
