package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/services"
	"github.com/tuantran0910/rainbow/pkg/headers"
	"github.com/tuantran0910/rainbow/pkg/utils/response"
)

type ProductController struct {
	productService services.IProductService
}

func NewProductController(productService services.IProductService) *ProductController {
	return &ProductController{
		productService: productService,
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
	if err != nil || page <= 0 {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid pagination's parameter page, page must be a positive integer").
			SetError(err.Error()).Respond(ctx)
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid pagination's parameter limit, limit must be a positive integer").
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
			SetMessage("Failed to fetch a list of products").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Parse output
	productResponses := make([]*dtos.GetProductResponse, 0)
	for _, product := range products {
		productResponses = append(productResponses, &dtos.GetProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
			DeletedAt: product.DeletedAt,
		})
	}
	data := &dtos.ListProductsResponse{
		Products: productResponses,
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

// GetProductById godoc
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
func (pc *ProductController) GetProductById(ctx *gin.Context) {
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
	product, err := pc.productService.GetProductByID(reqCtx, productId)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch a product").
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
	data := &dtos.GetProductResponse{
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
//	@Description	Create a product
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			req	body		dtos.CreateProductRequest	true	"Create Product Request"
//	@Success		204	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/products/{id} [post]
func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get the request body
	var productRequest dtos.CreateProductRequest
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
//	@Summary		Update Product by ID
//	@Description	Update a product by its ID
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			req	body		dtos.UpdateProductRequest	true	"Update Product Request"
//	@Success		204	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/products/{id} [patch]
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
	var productRequest dtos.UpdateProductRequest
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
//	@Summary		Delete Product by ID
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
		Respond(ctx)
}
