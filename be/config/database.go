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
	{Id: 2, Title: "Asset Manager", Slug: "assetManager", Description: "Manages assets", Activated: true, Created_at: time.Now()},
	{Id: 3, Title: "Department Head", Slug: "departmentHead", Description: "Oversees department", Activated: true, Created_at: time.Now()},
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
	{RoleId: 3, PermissionId: 6, AccessLevel: "can-request", Created_at: time.Now()},
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

func Int64Ptr(i int64) *int64 {
	return &i
}

var users = []entity.Users{{FirstName: "Admin",
	LastName: "Admin",
	RoleId:   1,
	Email:    "admin@gmail.com",
	Password: "$2a$10$uD2Sp/ceVMQs.Fxa9883Lejcy4QSiEsWFIihuosOkCqwQaCrs011.",
	IsActive: true},
	{FirstName: "Manager",
		LastName:       "asset 1",
		RoleId:         2,
		Email:          "ManagerAsset1@gmail.com",
		Password:       "$2a$10$Rkga1eAiQ4xSFSfIA.ZFyuraVz8lAE7/d.OsrVHb8Cd2J/KoVnkWu",
		IsActive:       true,
		DepartmentId:   Int64Ptr(1),
		IsAssetManager: true},
	{FirstName: "Manager",
		LastName:       "asset 2",
		RoleId:         2,
		Email:          "ManagerAsset2@gmail.com",
		Password:       "$2a$10$AGvvpScnwlpreNybde2RYOu3YwXWR5upqH4CYgY4kyrR9IUOS/2SC",
		IsActive:       true,
		DepartmentId:   Int64Ptr(2),
		IsAssetManager: true},
	{FirstName: "Manager",
		LastName:       "asset 3",
		RoleId:         2,
		Email:          "ManagerAsset3@gmail.com",
		Password:       "$2a$10$gPgRynYgAnJga.yDxY/E7OcjJFMFv4fsB3lL4lvnsvmpigYNMNJ2W",
		IsActive:       true,
		DepartmentId:   Int64Ptr(3),
		IsAssetManager: true},
	{FirstName: "Manager",
		LastName:       "asset 4",
		RoleId:         2,
		Email:          "ManagerAsset4@gmail.com",
		Password:       "$2a$10$Uu4bpMgDh5BqgCoxNNMD6ePiPXYJHOdCmDGf9JO7LflS6rxVo29t6",
		IsActive:       true,
		DepartmentId:   Int64Ptr(4),
		IsAssetManager: true},
	{FirstName: "Head",
		LastName:         "Department 1",
		RoleId:           3,
		Email:            "HeadDepartment1@gmail.com",
		Password:         "$2a$10$6FvD1M1qmg77Us8nURZJLu/aAL5b.nmazfg5vWTjF7uyUt4ysnVRC",
		IsActive:         true,
		DepartmentId:     Int64Ptr(1),
		IsHeadDepartment: true},
	{FirstName: "Head",
		LastName:         "Department 2",
		RoleId:           3,
		Email:            "HeadDepartment2@gmail.com",
		Password:         "$2a$10$5FwA7C7XsyrWP0hW5CXtfOx8bpLT8aGso6SZbIlMtdVNvvJA/KxWe",
		IsActive:         true,
		DepartmentId:     Int64Ptr(2),
		IsHeadDepartment: true},
	{FirstName: "Head",
		LastName:         "Department 3",
		RoleId:           3,
		Email:            "HeadDepartment3@gmail.com",
		Password:         "$2a$10$swIr.PfjoxEUxzuj/wDB9.5CyrLJ9ZHlV8x.H7H1ZUMaJCpy69Kha",
		IsActive:         true,
		DepartmentId:     Int64Ptr(3),
		IsHeadDepartment: true},
	{FirstName: "Head",
		LastName:         "Department 4",
		RoleId:           3,
		Email:            "HeadDepartment4@gmail.com",
		Password:         "$2a$10$S55fTYiB9AD4bYbyK47i8uiEneajsduoGd8eTsTT0tybfLhasGjzu",
		IsActive:         true,
		DepartmentId:     Int64Ptr(4),
		IsHeadDepartment: true},
	{FirstName: "Viewer",
		LastName:     "Department 1",
		RoleId:       4,
		Email:        "ViewerDepartment1@gmail.com",
		Password:     "$2a$10$HpKZlAE1EgXm2qSVUzDNY.Jl21nJdJoJF9N8Eo2h07WrFpKgd3hE6",
		IsActive:     true,
		DepartmentId: Int64Ptr(1),
	},
	{FirstName: "Viewer",
		LastName:     "Department 2",
		RoleId:       4,
		Email:        "ViewerDepartment2@gmail.com",
		Password:     "$2a$10$FbSLfcYefGmoqUFZWxIF2.TPb3ujjSsHCKhSYMP86VpEYozx6JCr6",
		IsActive:     true,
		DepartmentId: Int64Ptr(2),
	},
	{FirstName: "Viewer",
		LastName:     "Department 3",
		RoleId:       4,
		Email:        "ViewerDepartment3@gmail.com",
		Password:     "$2a$10$wVSYY0LmSYRXbEO3JyRWMu.JnNk.tsjCJgAMMSuWkm58eNMe2XmdW",
		IsActive:     true,
		DepartmentId: Int64Ptr(3),
	},
	{FirstName: "Viewer",
		LastName:     "Department 4",
		RoleId:       4,
		Email:        "ViewerDepartment4@gmail.com",
		Password:     "$2a$10$5t/A3R/jOLUxA2EFuCS/oeZA27i2YZ4PLBQAQ8/CK456dYUpMRrCa",
		IsActive:     true,
		DepartmentId: Int64Ptr(4),
	},
}

var locations = []entity.Locations{
	{Id: 1, LocationName: "307/12 Nguyen Van Troi,Ward 1, Tan Binh District,HCMC, Viet Nam"},
}
var departments = []entity.Departments{
	{Id: 1, LocationId: 1, DepartmentName: "PG1"},
	{Id: 2, LocationId: 1, DepartmentName: "PG2"},
	{Id: 3, LocationId: 1, DepartmentName: "PG3"},
	{Id: 4, LocationId: 1, DepartmentName: "PG4"},
}

func ConnectToDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(DB_DNS), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database. Error:", err)
	}
	createEnumSQL := `
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'asset_status') THEN
			CREATE TYPE asset_status AS ENUM (
				'New', 
				'In Use', 
				'Under Maintenance', 
				'Retired', 
				'Disposed'
			);
		END IF;
	END
	$$;
	`
	db.Exec(createEnumSQL)
	err = db.AutoMigrate(&entity.Roles{}, &entity.Permission{}, &entity.RolePermission{}, &entity.Users{}, &entity.UsersSessions{}, &entity.UserRbac{}, &entity.Locations{}, &entity.Departments{}, &entity.Categories{}, &entity.Assets{}, &entity.AssetLog{}, &entity.Assignments{}, &entity.RequestTransfer{}, &entity.Notifications{}, &entity.MaintenanceSchedules{}, &entity.MaintenanceNotifications{})
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
	for _, location := range locations {
		var existing entity.Locations
		db.Where("location_name = ?", existing.LocationName).FirstOrCreate(&existing, location)
	}
	for _, department := range departments {
		var existing entity.Departments
		db.Where("department_name = ?", existing.DepartmentName).FirstOrCreate(&existing, department)
	}
	for _, user := range users {
		var existing entity.Users
		db.Where("email = ?", user.Email).FirstOrCreate(&existing, user)
	}
	return db
}
