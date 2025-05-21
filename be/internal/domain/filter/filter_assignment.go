package filter

import (
	"fmt"

	"gorm.io/gorm"
)

type AssignmentFilter struct {
	EmailAssigned *string `form:"emailAssigned" json:"emailAssigned"`
	EmailAssign   *string `form:"emailAssign" json:"emailAssign"`
	AssetName     *string `form:"status" json:"status"`
	Page          int     `form:"page" json:"page"`
	Limit         int     `form:"limit" json:"limit"`
}

func (f *AssignmentFilter) ApplyFilter(db *gorm.DB, userId int64) *gorm.DB {
	db = db.Joins("join users as assigned_users  on assigned_users.id = assignments.user_id").
		Joins("join users as assigner_users on assigner_users.id = assignments.assign_by").
		Joins("join assets on assets.id = assignments.asset_id")
	if f.EmailAssigned != nil {
		str := fmt.Sprintf("%v", *f.EmailAssigned)
		str += "%"
		db = db.Where("LOWER(assigned_users.email) LIKE ?", str)
	}
	if f.EmailAssign != nil {
		str := fmt.Sprintf("%v", *f.EmailAssign)
		str += "%"
		db = db.Where("LOWER(assigner_users.email) LIKE ?", str)
	}
	if f.AssetName != nil {
		str := fmt.Sprintf("%v", *f.AssetName)
		str += "%"
		db = db.Where("LOWER(name) LIKE LOWER(?)", str)
	}
	return db.Preload("UserAssigned").Preload("UserAssign").Preload("Asset")
}
