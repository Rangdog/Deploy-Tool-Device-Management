package repository

import (
	"BE_Manage_device/internal/domain/entity"

	"gorm.io/gorm"
)

type PostgreSQLCompanyRepository struct {
	db *gorm.DB
}

func NewPostgreSQLCompanyRepository(db *gorm.DB) CompanyRepository {
	return &PostgreSQLCompanyRepository{db: db}
}

func (r *PostgreSQLCompanyRepository) Create(Company *entity.Company) (*entity.Company, error) {
	result := r.db.Create(Company)
	return Company, result.Error
}

func (r *PostgreSQLCompanyRepository) GetCompanyById(id int64) (*entity.Company, error) {
	var Company entity.Company
	result := r.db.Model(entity.Company{}).Where("id = ?", id).First(&Company)
	return &Company, result.Error
}

func (r *PostgreSQLCompanyRepository) GetCompanyBySuffixEmail(email string) (*entity.Company, error) {
	var Company entity.Company
	result := r.db.Model(entity.Company{}).Where("email = ?", email).First(&Company)
	return &Company, result.Error
}

func (r *PostgreSQLCompanyRepository) GetAllCompany() ([]*entity.Company, error) {
	var company []*entity.Company
	result := r.db.Model(entity.Company{}).Find(&company)
	return company, result.Error
}
