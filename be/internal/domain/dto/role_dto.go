package dto

type RoleResponse struct {
	Id          int64  `json:"id"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}
