package customer

import (
	"net/http"
	"strconv"
	"test-go/common"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository Repository
	Service    Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

// @Tags Customers
// @Summary Create a new customer
// @Description Create a new customer with the input payload
// @Accept  json
// @Produce  json
// @Param customer body CustomerCreateBody true "Customer Info"
// @Success 201 {object} map[string]interface{} "Created"
// @Failure 400 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /customers/ [post]
func (h *Handler) Create(c *gin.Context) {

	var body CustomerCreateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.ResponseError{Error: err.Error()})
		return
	}

	user := "email@mock.com" // Decode token at middleware

	input := &CustomerServiceCreateInput{
		CustomerCreateBody: body,
		CreatedBy:          user,
	}

	customerId, err := h.Service.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"customerId": customerId,
		},
	})
}

// @Tags Customers
// @Summary Get all customers
// @Description Retrieve a list of all customers with pagination and search keyword
// @Produce  json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param perPage query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param keyword query string false "Search keyword"
// @Success 200 {object} common.PaginatedResponse[CustomerTransformIndexOutput]
// @Failure 400 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /customers/ [get]
func (h *Handler) Index(c *gin.Context) {
	var query CustomerIndexQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, common.ResponseError{Error: err.Error()})
		return
	}

	customers, err := h.Service.FindAllAndCount(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{Error: err.Error()})
		return
	}

	var transformedCustomers []CustomerTransformIndexOutput
	for _, customer := range customers.Data {
		transformed := h.Service.TransformCustomerIndex(&customer)
		transformedCustomers = append(transformedCustomers, transformed)
	}

	pgQuery := common.PaginationQuery{
		Page:    query.Page,
		PerPage: query.PerPage,
	}

	c.JSON(http.StatusOK, common.BuildPaginatedResponseFromQuery(transformedCustomers, int(customers.TotalItems), pgQuery))
}

// @Tags Customers
// @Summary Get a customer by ID
// @Description Retrieve a single customer by their ID
// @Produce  json
// @Param id path int true "Customer ID"
// @Success 200 {object} CustomerShowResponse
// @Failure 400 {object} common.ResponseError
// @Failure 404 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /customers/{id} [get]
func (h *Handler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ResponseError{Error: "invalid id"})
		return
	}
	customer, err := h.Service.FindById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{Error: err.Error()})
		return
	}

	if customer == nil {
		c.JSON(http.StatusNotFound, common.ResponseError{Error: "customer not found"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": CustomerShowResponse{
			Id:        int(customer.Id),
			NameTh:    customer.NameTh,
			NameEn:    customer.NameEn,
			Email:     customer.Email,
			CreatedAt: customer.CreatedAt,
			CreatedBy: customer.CreatedBy,
		},
	})
}

// @Tags Customers
// @Summary Update a customer
// @Description Update an existing customer by ID with the input payload
// @Accept  json
// @Produce  json
// @Param id path uint true "Customer ID"
// @Param customer body CustomerUpdateBody true "Customer Info to update"
// @Success 200 {object} map[string]interface{} "Updated customer ID"
// @Failure 400 {object} common.ResponseError
// @Failure 404 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /customers/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	var body CustomerUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.ResponseError{Error: err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ResponseError{Error: "invalid customer ID"})
		return
	}

	existingCustomer, err := h.Service.FindById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{Error: err.Error()})
		return
	}
	if existingCustomer == nil {
		c.JSON(http.StatusNotFound, common.ResponseError{Error: "customer not found"})
		return
	}

	user := "email@mock.com" // Decode token at middleware

	input := &CustomerServiceUpdateInput{
		CustomerCreateBody: CustomerCreateBody{
			NameTh: body.NameTh,
			NameEn: body.NameEn,
			Email:  body.Email,
		},
		UpdatedBy: user,
	}

	customerId, err := h.Service.UpdateById(uint(id), input)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"customerId": customerId,
		},
	})
}

// @Tags Customers
// @Summary Delete a customer
// @Description Delete a customer by their ID
// @Param id path int true "Customer ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /customers/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ResponseError{Error: "invalid customer ID"})
		return
	}

	existingCustomer, err := h.Service.FindById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{Error: err.Error()})
		return
	}
	if existingCustomer == nil {
		c.JSON(http.StatusNotFound, common.ResponseError{Error: "customer not found"})
		return
	}

	if err := h.Service.DeleteById(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	customers := rg.Group("/customers")
	customers.POST("/", IsEmailExisted(h.Repository), h.Create)
	customers.GET("/", h.Index)
	customers.GET("/:id", h.Show)
	customers.PUT("/:id", IsEmailExisted(h.Repository), h.Update)
	customers.DELETE("/:id", h.Delete)
}
