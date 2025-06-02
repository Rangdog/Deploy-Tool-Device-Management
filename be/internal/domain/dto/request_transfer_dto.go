package dto

type CreateRequestTransferRequest struct {
	CategoryId  int64  `json:"categoryId" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type RequestTransferResponse struct {
	Id          int64                             `json:"id"`
	Status      string                            `json:"status"`
	User        UserResponseInRequestTransfer     `json:"user"`
	Category    CategoryResponseInRequestTransfer `json:"category"`
	Description string                            `json:"description"`
}

type UserResponseInRequestTransfer struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	DepartmentId int64  `json:"departmentId"`
}

type CategoryResponseInRequestTransfer struct {
	Id           int64  `json:"id"`
	CategoryName string `json:"categoryName"`
}

type ConfirmRequestTransferRequest struct {
	AssetId int64 `json:"assetId" binding:"required"`
}
