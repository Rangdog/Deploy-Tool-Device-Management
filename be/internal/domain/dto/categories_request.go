package dto

type CreateCategoryRequest struct {
	CategoryName string `json:"categoryName" binding:"required"`
}
