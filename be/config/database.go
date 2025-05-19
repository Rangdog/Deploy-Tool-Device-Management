package config

import (
	"BE_Manage_device/internal/domain/entity"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var roles = []entity.Roles{
	{Id: 1, Title: "Admin", Slug: "admin", Description: "Full access to system", Activated: true, Created_at: time.Now()},
	{Id: 2, Title: "Asset Manager", Slug: "asset-manager", Description: "Manages assets", Activated: true, Created_at: time.Now()},
	{Id: 3, Title: "Department Head", Slug: "department-head", Description: "Oversees department", Activated: true, Created_at: time.Now()},
	{Id: 4, Title: "Viewer", Slug: "viewer", Description: "Read-only access", Activated: true, Created_at: time.Now()},
}

var permissions = []entity.Permission{
	{Id: 1, Title: "User Management", Slug: "user-management", Description: "Manage system users", Activated: true, Created_at: time.Now()},
	{Id: 2, Title: "Role Assignment", Slug: "role-assignment", Description: "Assign roles to users", Activated: true, Created_at: time.Now()},
	{Id: 3, Title: "Create/Edit/Delete Assets", Slug: "manage-assets", Description: "Modify asset records", Activated: true, Created_at: time.Now()},
	{Id: 4, Title: "View Assets", Slug: "view-assets", Description: "View assets", Activated: true, Created_at: time.Now()},
	{Id: 5, Title: "Assign Assets to Users/Departments", Slug: "assign-assets", Description: "Assign assets", Activated: true, Created_at: time.Now()},
	{Id: 6, Title: "Transfer Assets Between Departments", Slug: "transfer-assets", Description: "Department-to-department transfer", Activated: true, Created_at: time.Now()},
	{Id: 7, Title: "Schedule Maintenance / Update Logs", Slug: "maintenance-logs", Description: "Log asset maintenance", Activated: true, Created_at: time.Now()},
	{Id: 8, Title: "Update Asset Lifecycle Stage", Slug: "lifecycle-update", Description: "Update lifecycle", Activated: true, Created_at: time.Now()},
	{Id: 9, Title: "Depreciation Management", Slug: "depreciation", Description: "Manage depreciation", Activated: true, Created_at: time.Now()},
	{Id: 10, Title: "Generate/View QR or Barcodes", Slug: "qr-barcodes", Description: "QR/barcode management", Activated: true, Created_at: time.Now()},
	{Id: 11, Title: "File Uploads", Slug: "file-uploads", Description: "Upload files to assets", Activated: true, Created_at: time.Now()},
	{Id: 12, Title: "View Dashboards / Reports", Slug: "dashboards", Description: "Access reports", Activated: true, Created_at: time.Now()},
	{Id: 13, Title: "Export Reports", Slug: "export-reports", Description: "Export data", Activated: true, Created_at: time.Now()},
	{Id: 14, Title: "Access Audit Logs / History", Slug: "audit-logs", Description: "Audit log access", Activated: true, Created_at: time.Now()},
	{Id: 15, Title: "Manage Categories, Departments, Locations", Slug: "manage-taxonomy", Description: "System categorization", Activated: true, Created_at: time.Now()},
	{Id: 16, Title: "Configure System Settings", Slug: "system-settings", Description: "General configuration", Activated: true, Created_at: time.Now()},
	{Id: 17, Title: "Integration Management", Slug: "integrations", Description: "ERP/API integration", Activated: true, Created_at: time.Now()},
	{Id: 18, Title: "Email Notifications / Alerts", Slug: "notifications", Description: "Email alert access", Activated: true, Created_at: time.Now()},
	{Id: 19, Title: "Access Billing / Subscriptions", Slug: "billing", Description: "Subscription management", Activated: true, Created_at: time.Now()},
}

var rolePermissions = []entity.RolePermission{
	// Admin - full access
	{RoleId: 1, PermissionId: 1, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 2, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 3, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 4, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 5, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 6, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 7, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 8, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 9, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 10, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 11, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 12, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 13, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 14, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 15, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 16, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 17, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 18, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 1, PermissionId: 19, AccessLevel: "full", Created_at: time.Now()},

	// Asset Manager
	{RoleId: 2, PermissionId: 3, AccessLevel: "limited", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 4, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 5, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 6, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 7, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 8, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 9, AccessLevel: "view", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 10, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 11, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 12, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 13, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 14, AccessLevel: "partial", Created_at: time.Now()},
	{RoleId: 2, PermissionId: 18, AccessLevel: "action", Created_at: time.Now()},

	// Department Head
	{RoleId: 3, PermissionId: 4, AccessLevel: "scoped", Created_at: time.Now()},
	{RoleId: 3, PermissionId: 7, AccessLevel: "view", Created_at: time.Now()},
	{RoleId: 3, PermissionId: 10, AccessLevel: "view", Created_at: time.Now()},
	{RoleId: 3, PermissionId: 12, AccessLevel: "scoped", Created_at: time.Now()},
	{RoleId: 3, PermissionId: 13, AccessLevel: "scoped", Created_at: time.Now()},
	{RoleId: 3, PermissionId: 18, AccessLevel: "receive", Created_at: time.Now()},

	// Viewer
	{RoleId: 4, PermissionId: 4, AccessLevel: "full", Created_at: time.Now()},
	{RoleId: 4, PermissionId: 10, AccessLevel: "scan", Created_at: time.Now()},
	{RoleId: 4, PermissionId: 12, AccessLevel: "scoped", Created_at: time.Now()},
	{RoleId: 4, PermissionId: 13, AccessLevel: "conditional", Created_at: time.Now()},
}

func ConnectToDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(DB_DNS), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database. Error:", err)
	}
	err = db.AutoMigrate(&entity.Roles{}, &entity.Permission{}, &entity.RolePermission{}, &entity.Users{}, &entity.UsersSessions{}, &entity.UserRbac{}, &entity.Locations{}, &entity.Departments{}, &entity.Categories{}, &entity.Assets{}, &entity.AssetLog{}, entity.Assignments{})
	if err != nil {
		log.Fatal("Error migrate to database. Error:", err)
	}

	for _, role := range roles {
		var existing entity.Roles
		db.Where("slug = ?", role.Slug).FirstOrCreate(&existing, role)
	}

	for _, permission := range permissions {
		var existing entity.Permission
		db.Where("slug = ?", permission.Slug).FirstOrCreate(&existing, permission)
	}

	for _, rolePermission := range rolePermissions {
		var existing entity.RolePermission
		db.Where("role_id = ? and permission_id = ?", rolePermission.RoleId, rolePermission.PermissionId).FirstOrCreate(&existing, rolePermission)
	}

	return db
}
