package dto

type CreateDepartmentRequest struct {
	DepartmentName string `json:"department_name"`
	LocationId     int64  `json:"location_id"`
}
