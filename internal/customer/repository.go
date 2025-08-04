package customer

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(customer *Customer) error
	FindAllAndCount(keyword string, page, perPage int) (CustomerServiceFindAllAndCount, error)
	FindById(id uint) (*Customer, error)
	UpdateById(customer *Customer) error
	DeleteById(id uint) error
	FindByEmail(email string, excludeId *uint) (*Customer, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(customer *Customer) error {
	return r.db.Create(customer).Error
}

func (r *repository) FindAllAndCount(keyword string, page, perPage int) (CustomerServiceFindAllAndCount, error) {
	var result CustomerServiceFindAllAndCount
	var customers []Customer
	var total int64

	db := r.db.Model(&Customer{})

	// กรองข้อมูลที่ยังไม่ถูกลบ (is_deleted = false หรือ IS NULL)
	db = db.Where("is_deleted IS NULL OR is_deleted = ?", false)

	// กรองด้วย keyword ถ้ามี
	if keyword != "" {
		likePattern := "%" + keyword + "%"
		db = db.Where("name_th ILIKE ? OR name_en ILIKE ? OR email ILIKE ?", likePattern, likePattern, likePattern)
	}

	// นับทั้งหมด
	if err := db.Count(&total).Error; err != nil {
		return result, err
	}

	offset := (page - 1) * perPage

	// ดึงข้อมูลตาม pagination
	if err := db.Limit(perPage).Offset(offset).Find(&customers).Error; err != nil {
		return result, err
	}

	result.Data = customers
	result.TotalItems = total

	return result, nil
}

func (r *repository) FindById(id uint) (*Customer, error) {
	var customer Customer
	err := r.db.
		Where("id = ? AND (is_deleted IS NULL OR is_deleted = false)", id).
		First(&customer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &customer, err
}

func (r *repository) UpdateById(customer *Customer) error {
	return r.db.Save(customer).Error
}

func (r *repository) DeleteById(id uint) error {
	return r.db.Model(&Customer{}).
		Where("id = ?", id).
		Updates(Customer{
			IsDeleted: true,
		}).Error
}

func (r *repository) FindByEmail(email string, excludeId *uint) (*Customer, error) {
	var customer Customer

	db := r.db.Model(&Customer{}).
		Where("email = ? AND (is_deleted IS NULL OR is_deleted = false)", email)

	if excludeId != nil {
		db = db.Where("id != ?", *excludeId)
	}

	err := db.First(&customer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &customer, nil
}
