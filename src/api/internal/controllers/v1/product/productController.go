package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/tuantran0910/rainbow/docs"
	model "github.com/tuantran0910/rainbow/internal/models/product"
	service "github.com/tuantran0910/rainbow/internal/services/product"
	"github.com/tuantran0910/rainbow/pkg/headers"
	"github.com/tuantran0910/rainbow/pkg/utils/response"
	"gorm.io/gorm"
)

type IProductController interface {
	GetProducts(ctx *gin.Context)
	GetProduct(ctx *gin.Context)
	CreateProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
}

type ProductController struct {
	productService service.IProductService
}

func NewProductController(db *gorm.DB) IProductController {
	return &ProductController{
		productService: service.NewProductService(db),
	}
}

// GetProducts godoc
//
//	@Summary		Get Products
//	@Description	Fetch a list of products
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"
//	@Param			limit	query		int	false	"Number of items per page"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/products [get]
func (pc *ProductController) GetProducts(ctx *gin.Context) {
	// Get pagination parameters from the query string
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid pagination's parameter page").
			SetError(err.Error()).Respond(ctx)
		return
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid pagination's parameter limit").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Get the context
	reqCtx := ctx.Request.Context()

	// Get a list of products
	products, pagination, err := pc.productService.GetProducts(reqCtx, page, limit)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot retrieve products").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Parse output
	data := make([]*model.ProductResponse, 0)
	for _, product := range products {
		data = append(data, &model.ProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
			DeletedAt: product.DeletedAt,
		})
	}

	// Set headers
	headers := headers.NewHeaders(data, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetPagination(pagination).
		SetMessage("Successfully retrieved products").
		SetData(data).
		Respond(ctx)
}

// GetProduct godoc
//
//	@Summary		Get Product by ID
//	@Description	Fetch a product by its ID
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/products/{id} [get]
func (pc *ProductController) GetProduct(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get ID from the URL
	productId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Get a product
	product, err := pc.productService.GetProduct(reqCtx, productId)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot retrieve product").
			SetError(err.Error()).Respond(ctx)
		return
	}

	if product == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("Product not found").
			Respond(ctx)
		return
	}

	// Parse output
	data := &model.ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
		DeletedAt: product.DeletedAt,
	}

	// Set headers
	headers := headers.NewHeaders(data, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved product").
		SetData(data).
		Respond(ctx)
}

// CreateProduct godoc
//
//	@Summary		Create Product
//	@Description	Create a new product and store it in the database
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			productRequest	body		model.ProductRequest	true	"Product Request"
//	@Success		201				{object}	response.APIResponse
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Router			/products [post]
func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get the request body
	var productRequest model.ProductRequest
	if err := ctx.ShouldBindJSON(&productRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Create a product
	if err := pc.productService.CreateProduct(reqCtx, productRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Cannot create product").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Set headers
	headers := headers.NewHeaders(nil, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusCreated).
		SetMessage("Successfully created product").
		Respond(ctx)
}

// UpdateProduct godoc
//
//	@Summary		Update Product
//	@Description	Update a product by its ID
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int						true	"Product ID"
//	@Param			productRequest	body		model.ProductRequest	true	"Product Request"
//	@Success		200				{object}	response.APIResponse
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Router			/products/{id} [put]
func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get ID from the URL
	productId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Get the request body
	var productRequest model.ProductRequest
	if err := ctx.ShouldBindJSON(&productRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Update a product
	if err := pc.productService.UpdateProduct(reqCtx, productId, productRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Cannot update product").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Set headers
	headers := headers.NewHeaders(nil, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully updated product").
		Respond(ctx)
}

// DeleteProduct godoc
//
//	@Summary		Delete Product
//	@Description	Delete a product by its ID
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		204	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/products/{id} [delete]
func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get ID from the URL
	productId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Delete a product
	if err := pc.productService.DeleteProduct(reqCtx, productId); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Cannot delete product").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Set headers
	headers := headers.NewHeaders(nil, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusNoContent).
		SetMessage("Successfully updated product").
		Respond(ctx)
}
