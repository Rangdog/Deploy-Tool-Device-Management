package dto

import "time"

type CreateMaintenanceSchedulesRequest struct {
	AssetId   int64     `json:"assetId" binding:"required"`
	StartDate time.Time `json:"startDate" binding:"required"`
	EndDate   time.Time `json:"endDate" binding:"required"`
}
