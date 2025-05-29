package dto

type AssignmentCreateRequest struct {
	UserId       *int64 `json:"userId"`
	AssetId      int64  `json:"assetId" binding:"required"`
	DepartmentId *int64 `json:"departmentId"`
}

type AssignmentUpdateRequest struct {
	UserId       *int64 `json:"userId"`
	DepartmentId *int64 `json:"departmentId"`
}

type AssignmentResponse struct {
	Id           int64                       `json:"id"`
	UserAssigned UsersAssignmentResponse     `json:"userAssigned"`
	UserAssign   UsersAssignmentResponse     `json:"userAssign"`
	Asset        UserAssignmentAssetResponse `json:"asset"`
	Department   DepartmentResponse          `json:"department"`
}

type UsersAssignmentResponse struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UserAssignmentAssetResponse struct {
	Id             int64  `json:"id"`
	AssetName      string `json:"assetName"`
	Status         string `json:"status"`
	FileAttachment string `json:"fileAttachment"`
	ImageUpload    string `json:"imageUpload"`
}
