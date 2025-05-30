package dto

type CreateRequestTransferRequest struct {
	AssetId      int64 `json:"assetId" binding:"required"`
	DepartmentId int64 `json:"departmentId" binding:"required"`
}

type RequestTransferResponse struct {
	Id         int64                          `json:"id"`
	Status     string                         `json:"status"`
	User       UserResponseInRequestTransfer  `json:"user"`
	Asset      AssetResponseInRequestTransfer `json:"asset"`
	Department DepartmentInRequestTransfer    `json:"Department"`
}

type UserResponseInRequestTransfer struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type AssetResponseInRequestTransfer struct {
	Id             int64  `json:"id"`
	AssetName      string `json:"assetName"`
	SerialNumber   string `json:"serialNumber"`
	ImageUpload    string `json:"image"`
	FileAttachment string `json:"file"`
	QrUrl          string `json:"qrUrl"`
}

type DepartmentInRequestTransfer struct {
	Id             int64            `json:"id"`
	DepartmentName string           `json:"departmentName"`
	Location       LocationResponse `json:"location"`
}
