package repository

import "BE_Manage_device/internal/domain/entity"

type RoleRepository interface {
	GetAllUserByRoleId(roleId int64) []*entity.Users
	GetAllUserByRoleSlug(slug string) []*entity.Users
	GetSlugByRoleId(id int64) string
	GetRoleBySlug(roleSlug string) *entity.Roles
	GetAllRole() []*entity.Roles
}
