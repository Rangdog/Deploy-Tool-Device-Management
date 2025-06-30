package repository

import (
	"BE_Manage_device/internal/domain/entity"
	"fmt"

	"gorm.io/gorm"
)

type PostgreSQLBillsRepository struct {
	db *gorm.DB
}

func NewPostgreSQLBillsRepository(db *gorm.DB) BillsRepository {
	return &PostgreSQLBillsRepository{db: db}
}

func (r *PostgreSQLBillsRepository) Create(bill *entity.Bill) (*entity.Bill, error) {
	var billNumber int64
	err := r.db.Raw("SELECT nextval('bill_number_seq')").Scan(&billNumber).Error
	if err != nil {
		return nil, err
	}
	bill.BillNumber = fmt.Sprintf("BILL-%08d", billNumber)
	result := r.db.Create(bill)
	return bill, result.Error
}

func (r *PostgreSQLBillsRepository) GetByBillNumber(billNumber string) (*entity.Bill, error) {
	var bill entity.Bill
	result := r.db.Model(entity.Bill{}).Where("bill_number =?", billNumber).Preload("CreateBy").Preload("Asset").Preload("CreateBy.Role").Preload("Asset.Category").Preload("Asset.Department").Preload("Asset.Department.Location").Preload("Asset.OnwerUser").First(&bill)
	return &bill, result.Error
}

func (r *PostgreSQLBillsRepository) GetDB() *gorm.DB {
	return r.db
}
