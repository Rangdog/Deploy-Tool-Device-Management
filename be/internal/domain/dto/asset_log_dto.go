package dto

type AssetLogsResponse struct {
	Action        string                         `json:"action"`
	Timestamp     string                         `json:"timeStamp"`
	ChangeSummary string                         `json:"changeSummary"`
	ByUser        UserResponseInAssetLog         `json:"byUserId"`
	AssignUser    *UserResponseInAssetLog        `json:"assignUserId"`
	Asset         AssetResponseInRequestTransfer `json:"asset"`
}

type UserResponseInAssetLog struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
