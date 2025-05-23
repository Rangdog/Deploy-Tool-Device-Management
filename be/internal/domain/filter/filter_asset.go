package filter

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type AssetFilter struct {
	AssetName    *string `form:"assetName" json:"assetName"`
	Status       *string `form:"status" json:"status"`
	CategoryId   *string `form:"categoryId" json:"categoryId"`
	Cost         *string `form:"cost" json:"cost"`
	SerialNumber *string `form:"serialNumber" json:"serialNumber"`
	Email        *string `form:"email" json:"email"`
	Page         int     `form:"page" json:"page"`
	Limit        int     `form:"limit" json:"limit"`
}

func (f *AssetFilter) ApplyFilter(db *gorm.DB, userId int64) *gorm.DB {
	db = db.Joins("JOIN user_rbacs on user_rbacs.asset_id = assets.id").
		Joins("JOIN roles on roles.id = user_rbacs.role_id").
		Joins("JOIN role_permissions on role_permissions.role_id = roles.id").
		Joins("JOIN permissions on permissions.id = role_permissions.permission_id").
		Joins("JOIN categories on categories.id = assets.category_id").
		Joins("JOIN users on users.id = assets.owner").
		Where("user_rbacs.user_id = ? and permissions.slug = ?", userId, "view-assets")
	if f.Status != nil {
		db = db.Where("status = ?", *f.Status)
	}
	if f.AssetName != nil {
		str := fmt.Sprintf("%v", *f.AssetName)
		str += "%"
		db = db.Where("LOWER(assets.asset_name) LIKE LOWER(?)", str)
	}
	if f.CategoryId != nil {
		parsedID, _ := strconv.ParseInt(*f.CategoryId, 10, 64)
		db = db.Where("categories.id = ?", parsedID)
	}
	if f.Cost != nil {
		if *f.Cost == "DESC" {
			db = db.Order("assets.cost DESC")
		} else if *f.Cost == "ASC" {
			db = db.Order("assets.cost ASC")
		} else {
			db = db.Order("assets.id ASC")
		}
	}
	if f.SerialNumber != nil {
		str := fmt.Sprintf("%v", *f.SerialNumber)
		str += "%"
		db = db.Where("LOWER(assets.serial_number) LIKE LOWER(?)", str)
	}
	if f.Email != nil {
		str := fmt.Sprintf("%v", *f.Email)
		str += "%"
		db = db.Where("LOWER(users.email) LIKE LOWER(?)", str)
	}
	db = db.Where("status != ?", "Disposed")
	return db.Preload("Category").Preload("Department").Preload("OnwerUser").Preload("Department.Location")
}
