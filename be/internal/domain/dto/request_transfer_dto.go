package dto

type CreateRequestTransferRequest struct {
	AssetId      int64 `json:"assetId" binding:"required"`
	DepartmentId int64 `json:"departmentId" binding:"required"`
}
