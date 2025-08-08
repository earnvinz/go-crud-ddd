package customer

import (
	"test-go/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	mockCreate          func(customer *Customer) error
	mockFindAllAndCount func(keyword string, page, perPage int) (CustomerServiceFindAllAndCount, error)
	mockFindById        func(id uint) (*Customer, error)
	mockUpdateById      func(customer *Customer) error
	mockDeleteById      func(id uint) error
	mockFindByEmail     func(email string, excludeId *uint) (*Customer, error)
}

func (m *mockRepository) Create(customer *Customer) error {
	if m.mockCreate != nil {
		return m.mockCreate(customer)
	}
	return nil
}

func (m *mockRepository) FindAllAndCount(keyword string, page, perPage int) (CustomerServiceFindAllAndCount, error) {
	if m.mockFindAllAndCount != nil {
		return m.mockFindAllAndCount(keyword, page, perPage)
	}
	return CustomerServiceFindAllAndCount{}, nil
}

func (m *mockRepository) FindById(id uint) (*Customer, error) {
	if m.mockFindById != nil {
		return m.mockFindById(id)
	}
	return nil, nil
}

func (m *mockRepository) UpdateById(customer *Customer) error {
	if m.mockUpdateById != nil {
		return m.mockUpdateById(customer)
	}
	return nil
}

func (m *mockRepository) DeleteById(id uint) error {
	if m.mockDeleteById != nil {
		return m.mockDeleteById(id)
	}
	return nil
}

func (m *mockRepository) FindByEmail(email string, excludeId *uint) (*Customer, error) {
	if m.mockFindByEmail != nil {
		return m.mockFindByEmail(email, excludeId)
	}
	return nil, nil
}

func TestService_Create(t *testing.T) {
	mockRepo := &mockRepository{
		mockCreate: func(c *Customer) error {
			c.Id = 123
			return nil
		},
	}

	svc := NewService(mockRepo)

	input := &CustomerServiceCreateInput{
		CustomerCreateBody: CustomerCreateBody{
			NameTh: "ทดสอบ",
			NameEn: "Test",
			Email:  "test@example.com",
		},
		CreatedBy: "unit@test.com",
	}

	id, err := svc.Create(input)
	assert.NoError(t, err)
	assert.Equal(t, uint(123), id)
}

func TestService_FindAllAndCount(t *testing.T) {
	mockRepo := &mockRepository{
		mockFindAllAndCount: func(keyword string, page, perPage int) (CustomerServiceFindAllAndCount, error) {
			return CustomerServiceFindAllAndCount{
				TotalItems: 1,
				Data: []Customer{
					{Id: 1, NameTh: "ทดสอบ", Email: "test@example.com"},
				},
			}, nil
		},
	}

	svc := NewService(mockRepo)

	keyword := "test"
	filter := CustomerIndexQuery{
		Keyword: &keyword,
		PaginationQuery: common.PaginationQuery{
			Page:    1,
			PerPage: 10,
		},
	}

	result, err := svc.FindAllAndCount(filter)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.TotalItems)

	assert.Len(t, result.Data, 1)
	assert.Equal(t, "test@example.com", result.Data[0].Email)
}

func TestService_UpdateById(t *testing.T) {
	updated := false
	mockRepo := &mockRepository{
		mockUpdateById: func(c *Customer) error {
			updated = true
			assert.Equal(t, uint(1), c.Id)
			assert.Equal(t, "Updated Name", c.NameTh)
			return nil
		},
	}

	svc := NewService(mockRepo)

	input := &CustomerServiceUpdateInput{
		CustomerCreateBody: CustomerCreateBody{
			NameTh: "Updated Name",
			NameEn: "Updated EN",
			Email:  "updated@example.com",
		},
		UpdatedBy: "unit@test.com",
	}

	id, err := svc.UpdateById(1, input)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)
	assert.True(t, updated)
}

func TestService_DeleteById(t *testing.T) {
	deleted := false
	mockRepo := &mockRepository{
		mockDeleteById: func(id uint) error {
			deleted = true
			assert.Equal(t, uint(1), id)
			return nil
		},
	}

	svc := NewService(mockRepo)

	err := svc.DeleteById(1)
	assert.NoError(t, err)
	assert.True(t, deleted)
}

func TestService_FindById(t *testing.T) {
	mockRepo := &mockRepository{
		mockFindById: func(id uint) (*Customer, error) {
			return &Customer{Id: id, Email: "findbyid@example.com"}, nil
		},
	}

	svc := NewService(mockRepo)

	cust, err := svc.FindById(1)
	assert.NoError(t, err)
	assert.NotNil(t, cust)
	assert.Equal(t, "findbyid@example.com", cust.Email)
}

func TestService_FindByEmail(t *testing.T) {
	mockRepo := &mockRepository{
		mockFindByEmail: func(email string, excludeId *uint) (*Customer, error) {
			if email == "exists@example.com" {
				return &Customer{Id: 1, Email: email}, nil
			}
			return nil, nil
		},
	}

	svc := NewService(mockRepo)

	cust, err := svc.FindByEmail("exists@example.com", nil)
	assert.NoError(t, err)
	assert.NotNil(t, cust)
	assert.Equal(t, "exists@example.com", cust.Email)

	cust, err = svc.FindByEmail("notfound@example.com", nil)
	assert.NoError(t, err)
	assert.Nil(t, cust)
}
