package repository

import (
	"BE_Manage_device/internal/domain/entity"
	"errors"

	"gorm.io/gorm"
)

type PostgreSQLUserRepository struct {
	db *gorm.DB
}

func NewPostgreSQLUserRepository(db *gorm.DB) UserRepository {
	return &PostgreSQLUserRepository{db: db}
}

func (r *PostgreSQLUserRepository) Create(users *entity.Users) error {
	if users.FirstName == "" {
		return errors.New("name can't not blank")
	}
	if users.LastName == "" {
		return errors.New("name can't not blank")
	}
	if users.Password == "" {
		return errors.New("password can't not blank")
	}
	if err := r.db.Create(users).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgreSQLUserRepository) FindByToken(token string) (*entity.Users, error) {
	var users = &entity.Users{}
	result := r.db.Model(&entity.Users{}).Where("token = ?", token).Find(users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) Update(users *entity.Users) error {
	result := r.db.Model(&entity.Users{}).Where("email = ?", users.Email).Update("is_active", true)
	return result.Error
}

func (r *PostgreSQLUserRepository) UpdatePassword(users *entity.Users) error {
	result := r.db.Model(&entity.Users{}).Where("email = ?", users.Email).Update("password", users.Password)
	return result.Error
}

func (r *PostgreSQLUserRepository) FindByEmail(email string) (*entity.Users, error) {
	users := &entity.Users{}
	result := r.db.Model(entity.Users{}).Where("email = ?", email).Preload("Role").Preload("Department").First(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) FindByUserId(userId int64) (*entity.Users, error) {
	users := &entity.Users{}
	result := r.db.Model(entity.Users{}).Where("id = ?", userId).Preload("Role").Preload("Department").First(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) DeleteUser(email string) error {
	result := r.db.Where("email = ?", email).Delete(&entity.Users{})
	return result.Error
}

func (r *PostgreSQLUserRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PostgreSQLUserRepository) GetAllUser(companyId int64) []*entity.Users {
	var users = []*entity.Users{}
	result := r.db.Model(entity.Users{}).Where("company_id = ?", companyId).Preload("Role").Preload("Department").Find(&users)
	if result.Error != nil {
		return nil
	}
	return users
}

func (r *PostgreSQLUserRepository) UpdateUser(user *entity.Users) (*entity.Users, error) {
	var userUpdate = entity.Users{}
	updates := map[string]interface{}{}
	if user.FirstName != "" {
		updates["first_name"] = user.FirstName
	}
	if user.LastName != "" {
		updates["last_name"] = user.LastName
	}
	if user.RoleId != 0 {
		updates["role_id"] = user.RoleId
	}
	if user.Avatar != "" {
		updates["avatar"] = user.Avatar
	}
	err := r.db.Model(&userUpdate).Where("id = ?", user.Id).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// Trả về bản ghi sau khi cập nhật (tuỳ bạn muốn lấy lại hay không)
	err = r.db.Preload("Role").First(&userUpdate, user.Id).Error
	if err != nil {
		return nil, err
	}

	return &userUpdate, nil
}

func (r *PostgreSQLUserRepository) GetUserAssetManageOfDepartment(departmentId int64) (*entity.Users, error) {
	user := entity.Users{}
	result := r.db.Model(entity.Users{}).Where("department_id = ? and is_asset_manager = true", departmentId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *PostgreSQLUserRepository) FindByEmailForLogin(email string) (*entity.Users, error) {
	users := &entity.Users{}
	result := r.db.Model(entity.Users{}).Where("email = ? and is_active = true", email).First(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) GetAllUserRoleEmployeeOfDepartment(departmentTd int64) ([]*entity.Users, error) {
	users := []*entity.Users{}
	result := r.db.Debug().Model(entity.Users{}).Joins("join roles on roles.id = users.role_id").Where("department_id = ?", departmentTd).Where("roles.slug = ?", "employee").Preload("Role").Preload("Department").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) GetAllUserRoleManagerOfDepartment(departmentTd int64) ([]*entity.Users, error) {
	users := []*entity.Users{}
	result := r.db.Debug().Model(entity.Users{}).Joins("join roles on roles.id = users.role_id").Where("department_id = ?", departmentTd).Where("roles.slug = ?", "assetManager").Preload("Role").Preload("Department").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) UpdateDepartment(userId int64, departmentId int64) error {
	result := r.db.Model(entity.Users{}).Where("id = ?", userId).Update("department_id", departmentId)
	return result.Error
}

func (r *PostgreSQLUserRepository) CheckHeadDep(depId int64) error {
	var count int64
	r.db.Model(&entity.Users{}).
		Where("department_id = ? AND is_head_department = ?", depId, true).
		Count(&count)
	if count > 0 {
		return errors.New("the department already has a head of department")
	}
	return nil
}

func (r *PostgreSQLUserRepository) CheckManagerDep(depId int64) error {
	var count int64
	r.db.Model(&entity.Users{}).
		Where("department_id = ? AND is_asset_manager = ?", depId, true).
		Count(&count)
	if count > 0 {
		return errors.New("the department already has a manager asset of department")
	}
	return nil
}

func (r *PostgreSQLUserRepository) UpdateHeadDep(id int64, isHead bool) error {
	result := r.db.Model(entity.Users{}).Where("id = ?", id).Update("is_head_department", isHead)
	return result.Error
}

func (r *PostgreSQLUserRepository) UpdateManagerDep(id int64, isManager bool) error {
	result := r.db.Model(entity.Users{}).Where("id = ?", id).Update("is_asset_manager", isManager)
	return result.Error
}

func (r *PostgreSQLUserRepository) UpdateCanExport(id int64, canExport bool) error {
	result := r.db.Model(entity.Users{}).Where("id = ?", id).Update("can_export", canExport)
	return result.Error
}

func (r *PostgreSQLUserRepository) GetUserNotHaveDep() ([]*entity.Users, error) {
	var user []*entity.Users
	result := r.db.Model(entity.Users{}).Joins("join roles on roles.id = users.role_id").Where("roles.slug != ?", "admin").Where("department_id is null").Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *PostgreSQLUserRepository) GetUserRoleAdmin() ([]*entity.Users, error) {
	var users []*entity.Users
	result := r.db.Model(entity.Users{}).Where("role_id = (select id from roles where slug = ?)", "admin").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) FindManager(userId int64) (*entity.Users, error) {
	var user entity.Users
	result := r.db.Model(entity.Users{}).Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.DepartmentId == nil {
		return nil, errors.New("user haven't department")
	}
	depId := user.DepartmentId
	var userManager entity.Users
	result = r.db.Model(entity.Users{}).Where("department_id = ? and is_asset_manager = true", depId).First(&userManager)
	return &userManager, result.Error
}
