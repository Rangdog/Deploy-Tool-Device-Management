package dto

type CreateLocationRequest struct {
	LocationName string `json:"locationName" binding:"required"`
}
