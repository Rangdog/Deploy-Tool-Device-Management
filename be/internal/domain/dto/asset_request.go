package dto

type AssetResponse struct {
	ID             int64              `json:"id"`
	AssetName      string             `json:"assetName"`
	PurchaseDate   string             `json:"purchaseDate"`
	Cost           float64            `json:"cost"`
	Owner          OwnerResponse      `json:"owner,omitempty"`
	WarrantExpiry  string             `json:"warrantExpiry"`
	Status         string             `json:"status"`
	SerialNumber   string             `json:"serialNumber"`
	FileAttachment string             `json:"fileAttachment"`
	ImageUpload    string             `json:"imageUpload"`
	Category       CategoryResponse   `json:"category"`
	QrURL          string             `json:"qrUrl"`
	Department     DepartmentResponse `json:"department"`
}

type CategoryResponse struct {
	ID           int64  `json:"id"`
	CategoryName string `json:"categoryName"`
}

type DepartmentResponse struct {
	ID             int64            `json:"id"`
	DepartmentName string           `json:"departmentName"`
	Location       LocationResponse `json:"location"`
}

type LocationResponse struct {
	ID           int64  `json:"id"`
	LocationName string `json:"locationAddress"`
}

type OwnerResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type DashboardSummary struct {
	TotalAssets      int `json:"total_assets"`
	Assigned         int `json:"assigned"`
	UnderMaintenance int `json:"under_maintenance"`
	Retired          int `json:"retired"`
}

type GetAssetsByCateOfDepartmentRequest struct {
	CategoryId   int64 `form:"categoryId"`
	DepartmentId int64 `form:"departmentId"`
}

type RetiredAssetRequest struct {
	ResidualValue float64 `json:"residualValue" binding:"required"`
}
