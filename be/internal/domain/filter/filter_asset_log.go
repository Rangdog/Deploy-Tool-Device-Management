package filter

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AssetLogFilter struct {
	Action    *string `form:"action" json:"action"`
	StartTime *string `form:"startTime" json:"startTime"`
	EndTime   *string `form:"endTime" json:"endTime"`
	DepId     *int64
}

func (f *AssetLogFilter) ApplyFilter(db *gorm.DB, assetId int64) *gorm.DB {
	db = db.Where("asset_id = ?", assetId)
	if f.Action != nil {
		str := fmt.Sprintf("%v", strings.ToLower(*f.Action))
		str += "%"
		db = db.Where("LOWER(action) LIKE ?", str)
	}
	if f.StartTime != nil && f.EndTime != nil {
		if t, err := time.Parse(time.RFC3339Nano, *f.StartTime); err == nil {
			db = db.Where("timestamp>=?", t)
		}
		if t, err := time.Parse(time.RFC3339Nano, *f.EndTime); err == nil {
			db = db.Where("timestamp<=?", t)
		}
	}
	if f.DepId != nil {
		db = db.Joins("join assets on assets.id = asset_logs.asset_id").Where("assets.department_id = ?", *f.DepId)
	}
	return db.Preload("ByUser").Preload("AssignUser").Preload("Asset").Order("id ASC")
}
