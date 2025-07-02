package filter

import "gorm.io/gorm"

type MonthlySummaryFilter struct {
	Month     int64 `form:"month"`
	Year      int64 `form:"year"`
	CompanyId int64 `form:"-"`
}

func (f *MonthlySummaryFilter) ApplyFilter(db *gorm.DB) *gorm.DB {
	db.Where("monthly_summaries.company_id = ?", f.CompanyId)
	db.Where("monthly_summaries.month = ? and monthly_summaries.year = ?", f.Month, f.Year)
	return db
}
