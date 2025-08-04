package common

import "math"

type PaginationQuery struct {
	Page    int `form:"page" binding:"min=1"`
	PerPage int `form:"perPage" binding:"min=1"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type PaginatedResponse[T any] struct {
	Data []T `json:"data"`
	PaginationMeta
}

func BuildPaginatedResponseFromQuery[T any](input []T, totalItems int, query PaginationQuery) PaginatedResponse[T] {
	page := query.Page
	if page < 1 {
		page = 1
	}
	perPage := query.PerPage
	if perPage < 1 {
		perPage = 10
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	if input == nil {
		input = []T{}
	}

	response := PaginatedResponse[T]{
		PaginationMeta: PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
		Data: input,
	}

	return response
}
