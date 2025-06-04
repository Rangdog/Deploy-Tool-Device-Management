package dto

import "time"

type CreateMaintenanceSchedulesRequest struct {
	AssetId   int64     `json:"assetId" binding:"required"`
	StartDate time.Time `json:"startDate" binding:"required"`
	EndDate   time.Time `json:"endDate" binding:"required"`
}

type UpdateMaintenanceSchedulesRequest struct {
	StartDate time.Time `json:"startDate" binding:"required"`
	EndDate   time.Time `json:"endDate" binding:"required"`
}

type MaintenanceSchedulesResponse struct {
	Id        int64                               `json:"id"`
	StartDate string                              `json:"startDate"`
	EndDate   string                              `json:"endDate"`
	Asset     AssetResponseInMaintenanceSchedules `json:"asset"`
}

type AssetResponseInMaintenanceSchedules struct {
	Id             int64  `json:"id"`
	AssetName      string `json:"assetName"`
	Status         string `json:"status"`
	FileAttachment string `json:"fileAttachment"`
	ImageUpload    string `json:"imageUpload"`
}
