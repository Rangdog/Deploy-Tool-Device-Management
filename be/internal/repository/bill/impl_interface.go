package repository

import (
	"BE_Manage_device/internal/domain/entity"
	"fmt"
	"time"

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

func monthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}

func (r *PostgreSQLBillsRepository) GetAllBillOfMonth(time time.Time, companyId int64) ([]*entity.Bill, error) {
	y := time.Year()
	m := time.Month()
	first, last := monthInterval(y, m)
	var bills []*entity.Bill
	result := r.db.Model(entity.Bill{}).Where("company_id = ?", companyId).Where("create_at >= ? and create_at <= ?", first, last).Preload("Asset").Preload("Asset.Category").Find(&bills)
	return bills, result.Error
}

func (r *PostgreSQLBillsRepository) GetAllBillUnpaid(companyId int64) ([]*entity.Bill, error) {
	var bills []*entity.Bill
	result := r.db.Model(entity.Bill{}).Where("company_id = ?", companyId).Where("status_bill = ?", "UnPaid").Preload("CreateBy").Preload("Asset").Preload("CreateBy.Role").Preload("Asset.Category").Preload("Asset.Department").Preload("Asset.Department.Location").Preload("Asset.OnwerUser").Find(&bills)
	return bills, result.Error
}

func (r *PostgreSQLBillsRepository) UpdatePaid(billNumberStr string) error {
	result := r.db.Model(entity.Bill{}).Where("bill_number = ?", billNumberStr).Update("status_bill", "Paid")
	return result.Error
}
