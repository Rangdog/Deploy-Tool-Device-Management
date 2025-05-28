package dto

type AssetLogsResponse struct {
	Action        string                  `json:"action"`
	Timestamp     string                  `json:"timeStamp"`
	ChangeSummary string                  `json:"changeSummary"`
	Users         UserResponseInAssetLog  `json:"user"`
	Department    DepartmentResponse      `json:"department"`
	Asset         AssetResponseInAssetLog `json:"asset"`
}

type UserResponseInAssetLog struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type AssetResponseInAssetLog struct {
	AssetName      string `json:"assetName"`
	SerialNumber   string `json:"serialNumber"`
	ImageUpload    string `json:"image"`
	FileAttachment string `json:"file"`
	QrUrl          string `json:"qrUrl"`
}
