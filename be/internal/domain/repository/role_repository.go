package repository

import "BE_Manage_device/internal/domain/entity"

type RoleRepository interface {
	GetAllUserByRoleId(roleId int64) []*entity.Users
	GetAllUserByRoleTitle(title string) []*entity.Users
	GetTitleByRoleId(id int64) string
	GetRoleByTile(roleTitle string) *entity.Roles
	GetAllRole() []*entity.Roles
}
