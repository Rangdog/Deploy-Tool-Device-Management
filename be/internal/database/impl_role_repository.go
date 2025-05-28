package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

	"gorm.io/gorm"
)

type PostgreSQLRoleRepository struct {
	db *gorm.DB
}

func NewPostgreSQLRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &PostgreSQLRoleRepository{db: db}
}

func (r *PostgreSQLRoleRepository) GetAllUserByRoleId(roleId int64) []*entity.Users {
	var users = []*entity.Users{}
	result := r.db.Model(entity.Roles{}).Joins("Join users on users.role_id = roles.id").Where("roles.id = ?", roleId).Find(users)
	if result.Error != nil {
		return nil
	}
	return users
}

func (r *PostgreSQLRoleRepository) GetAllUserByRoleSlug(slug string) []*entity.Users {
	var users = []*entity.Users{}
	result := r.db.Model(entity.Roles{}).Joins("Join users on users.role_id = roles.id").Where("roles.slug = ?", slug).Find(users)
	if result.Error != nil {
		return nil
	}
	return users
}

func (r *PostgreSQLRoleRepository) GetSlugByRoleId(id int64) string {
	role := entity.Roles{}
	r.db.Model(entity.Roles{}).Where("id = ?", id).First(&role)
	return role.Slug
}

func (r *PostgreSQLRoleRepository) GetRoleBySlug(slug string) *entity.Roles {
	roles := entity.Roles{}
	r.db.Model(entity.Roles{}).Where("slug = ?", slug).Find(&roles)
	return &roles
}

func (r *PostgreSQLRoleRepository) GetAllRole() []*entity.Roles {
	roles := []*entity.Roles{}
	r.db.Model(entity.Roles{}).Find(&roles)
	return roles
}
