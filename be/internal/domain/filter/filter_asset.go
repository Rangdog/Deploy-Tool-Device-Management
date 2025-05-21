package filter

import (
	"fmt"

	"gorm.io/gorm"
)

type AssetFilter struct {
	AssetName *string `form:"asset_name" json:"asset_name"`
	Status    *string `form:"status" json:"status"`
	Page      int     `form:"page" json:"page"`
	Limit     int     `form:"limit" json:"limit"`
}

func (f *AssetFilter) ApplyFilter(db *gorm.DB, userId int64) *gorm.DB {
	db = db.Joins("JOIN user_rbacs on user_rbacs.asset_id = assets.id").
		Joins("JOIN roles on roles.id = user_rbacs.role_id").
		Joins("JOIN role_permissions on role_permissions.role_id = roles.id").
		Joins("JOIN permissions on permissions.id = role_permissions.permission_id").
		Where("user_rbacs.user_id = ? and permissions.slug = ?", userId, "view-assets")
	if f.Status != nil {
		db = db.Where("status = ?", *f.Status)
	}
	if f.AssetName != nil {
		str := fmt.Sprintf("%v", *f.AssetName)
		str += "%"
		db = db.Where("LOWER(name) LIKE LOWER(?)", str)
	}
	return db.Preload("Category").Preload("Department").Preload("OnwerUser").Preload("Department.Location")
}
