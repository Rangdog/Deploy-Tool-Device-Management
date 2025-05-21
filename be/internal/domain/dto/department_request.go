package dto

type CreateDepartmentRequest struct {
	DepartmentName string `json:"departmentName" binding:"required"`
	LocationId     int64  `json:"locationId" binding:"required"`
}
