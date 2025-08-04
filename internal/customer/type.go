package customer

import (
	"test-go/common"
	"time"
)

type CustomerShowResponse struct {
	Id        int       `json:"id"`
	NameTh    string    `json:"nameTh"`
	NameEn    string    `json:"nameEn"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
}

type CustomerCreateBody struct {
	NameTh string `json:"nameTh" binding:"required" example:"สมชาย"`
	NameEn string `json:"nameEn" binding:"required" example:"Somchai"`
	Email  string `json:"email" binding:"required,email" example:"somchai@example.com"`
}

type CustomerServiceCreateInput struct {
	CustomerCreateBody
	CreatedBy string
}

type CustomerIndexQuery struct {
	common.PaginationQuery
	Keyword *string `form:"keyword" example:"search term"`
}

type CustomerTransformIndexOutput struct {
	Id        uint
	NameTh    string
	NameEn    string
	Email     string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}

type CustomerServiceFindAllAndCount struct {
	Data       []Customer
	TotalItems int64
}

type CustomerUpdateBody struct {
	CustomerCreateBody
}

type CustomerServiceUpdateInput struct {
	CustomerCreateBody
	UpdatedBy string
}
