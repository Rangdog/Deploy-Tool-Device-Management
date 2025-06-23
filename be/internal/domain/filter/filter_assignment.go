package filter

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type AssignmentFilter struct {
	EmailAssigned *string `form:"emailAssigned" json:"emailAssigned"`
	EmailAssign   *string `form:"emailAssign" json:"emailAssign"`
	AssetName     *string `form:"assetName" json:"assetName"`
}

func (f *AssignmentFilter) ApplyFilter(db *gorm.DB, userId int64) *gorm.DB {
	db = db.Joins("join users as assigned_users  on assigned_users.id = assignments.user_id").
		Joins("join users as assigner_users on assigner_users.id = assignments.assign_by").
		Joins("join assets on assets.id = assignments.asset_id")
	if f.EmailAssigned != nil {
		str := fmt.Sprintf("%v", strings.ToLower(*f.EmailAssigned))
		str += "%"
		db = db.Where("LOWER(assigned_users.email) LIKE ?", str)
	}
	if f.EmailAssign != nil {
		str := fmt.Sprintf("%v", strings.ToLower(*f.EmailAssign))
		str += "%"
		db = db.Where("LOWER(assigner_users.email) LIKE ?", str)
	}
	if f.AssetName != nil {
		str := "%" + fmt.Sprintf("%v", *f.AssetName)
		str += "%"
		db = db.Where("LOWER(assets.asset_name) LIKE LOWER(?)", str)
	}
	return db.Preload("UserAssigned").Preload("UserAssign").Preload("Asset").Preload("Department").Preload("Department.Location").Order("assignments.id ASC")
}
