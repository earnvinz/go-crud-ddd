package customer

import (
	"time"
)

type Service interface {
	Create(customer *CustomerServiceCreateInput) (uint, error)
	FindAllAndCount(filter CustomerIndexQuery) (CustomerServiceFindAllAndCount, error)
	UpdateById(id uint, input *CustomerServiceUpdateInput) (uint, error)
	DeleteById(id uint) error
	TransformCustomerIndex(customer *Customer) CustomerTransformIndexOutput
	FindById(id uint) (*Customer, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(input *CustomerServiceCreateInput) (uint, error) {
	now := time.Now()
	customer := &Customer{
		NameTh:    input.NameTh,
		NameEn:    input.NameEn,
		Email:     input.Email,
		CreatedBy: input.CreatedBy,
		CreatedAt: now,
		UpdatedBy: input.CreatedBy,
		UpdatedAt: now,
	}

	err := s.repo.Create(customer)
	if err != nil {
		return 0, err
	}
	return customer.Id, nil
}

func (s *service) FindAllAndCount(filter CustomerIndexQuery) (CustomerServiceFindAllAndCount, error) {
	keyword := ""
	if filter.Keyword != nil {
		keyword = *filter.Keyword
	}

	return s.repo.FindAllAndCount(keyword, filter.Page, filter.PerPage)
}

func (s *service) UpdateById(id uint, input *CustomerServiceUpdateInput) (uint, error) {
	now := time.Now()
	customer := &Customer{
		NameTh: input.NameTh,
		NameEn: input.NameEn,
		Email:  input.Email,

		UpdatedBy: input.UpdatedBy,
		UpdatedAt: now,
	}
	return customer.Id, nil
}

func (s *service) DeleteById(id uint) error {
	return s.repo.DeleteById(id)
}

func (s *service) TransformCustomerIndex(customer *Customer) CustomerTransformIndexOutput {
	return CustomerTransformIndexOutput{
		Id:        customer.Id,
		NameTh:    customer.NameTh,
		NameEn:    customer.NameEn,
		Email:     customer.Email,
		CreatedAt: customer.CreatedAt,
		CreatedBy: customer.CreatedBy,
		UpdatedAt: customer.UpdatedAt,
		UpdatedBy: customer.UpdatedBy,
	}
}

func (s *service) FindById(id uint) (*Customer, error) {
	return s.repo.FindById(id)
}
