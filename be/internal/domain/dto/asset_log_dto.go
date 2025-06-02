package dto

type AssetLogsResponse struct {
	Action        string                  `json:"action"`
	Timestamp     string                  `json:"timeStamp"`
	ChangeSummary string                  `json:"changeSummary"`
	ByUser        UserResponseInAssetLog  `json:"byUserId"`
	AssignUser    *UserResponseInAssetLog `json:"assignUserId"`
	Asset         AssetResponseInAssetLog `json:"asset"`
}

type UserResponseInAssetLog struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type AssetResponseInAssetLog struct {
	Id             int64  `json:"id"`
	AssetName      string `json:"assetName"`
	SerialNumber   string `json:"serialNumber"`
	ImageUpload    string `json:"image"`
	FileAttachment string `json:"file"`
	QrUrl          string `json:"qrUrl"`
}
