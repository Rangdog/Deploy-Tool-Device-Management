package repository

import "BE_Manage_device/internal/domain/entity"

type RoleRepository interface {
	GetAllUserByRoleId(roleId int64) []*entity.Users
	GetAllUserByRoleSlug(title string) []*entity.Users
	GetSlugByRoleId(id int64) string
	GetRoleBySlug(roleTitle string) *entity.Roles
	GetAllRole() []*entity.Roles
}
