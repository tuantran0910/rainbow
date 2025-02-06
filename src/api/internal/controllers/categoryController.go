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

type CategoryController struct {
	categoryService services.ICategoryService
}

func NewCategoryController(categoryService services.ICategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// GetCategories godoc
//
//	@Summary		Get a list of categories
//	@Description	Get a list of categories
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"
//	@Param			limit	query		int	false	"Number of items per page"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/categories [get]
func (cc *CategoryController) GetCategories(ctx *gin.Context) {
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

	reqCtx := ctx.Request.Context()
	categories, pagination, err := cc.categoryService.GetCategories(reqCtx, page, limit)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to get categories").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	categoryResponses := make([]*dtos.GetCategoryResponse, 0)
	for _, category := range categories {
		categoryResponses = append(categoryResponses, &dtos.GetCategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Slug:      category.Slug,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
			DeletedAt: category.DeletedAt,
		})
	}
	data := &dtos.ListCategoriesResponse{
		Categories: categoryResponses,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetPagination(pagination).
		SetMessage("Successfully retrieved categories").
		SetData(data).
		Respond(ctx)
}

// GetCategoryById godoc
//
//	@Summary		Get a category
//	@Description	Get a category by its ID
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Category ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/categories/{id} [get]
func (cc *CategoryController) GetCategoryById(ctx *gin.Context) {
	categoryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
		return
	}

	reqCtx := ctx.Request.Context()
	category, err := cc.categoryService.GetCategoryById(reqCtx, categoryId)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to get the category").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	if category == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("Category not found").
			Respond(ctx)
		return
	}

	data := &dtos.GetCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
		DeletedAt: category.DeletedAt,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved category").
		SetData(data).
		Respond(ctx)
}

// GetCategoryBySlug godoc
//
//	@Summary		Get a category
//	@Description	Get a category by its slug
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"Category slug"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/categories/{slug} [get]
func (cc *CategoryController) GetCategoryBySlug(ctx *gin.Context) {
	categorySlug := ctx.Param("slug")

	reqCtx := ctx.Request.Context()
	category, err := cc.categoryService.GetCategoryBySlug(reqCtx, categorySlug)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to get the category").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	if category == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("Category not found").
			Respond(ctx)
		return
	}

	data := &dtos.GetCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
		DeletedAt: category.DeletedAt,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved category").
		SetData(data).
		Respond(ctx)
}

// CreateCategory godoc
//
//	@Summary		Create a category
//	@Description	Create a category
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.CreateCategoryRequest	true	"Category information"
//	@Success		201		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/categories [post]
func (cc *CategoryController) CreateCategory(ctx *gin.Context) {
	var categoryRequest dtos.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&categoryRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid request body").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	currentUserId, ok := ctx.Get("user_id")
	if !ok {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot get the current user's ID").
			Respond(ctx)
		return
	}

	reqCtx := ctx.Request.Context()
	if err := cc.categoryService.CreateCategory(reqCtx, categoryRequest, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to create the category").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	headers := headers.NewHeaders(nil, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusCreated).
		SetMessage("Successfully created the category").
		Respond(ctx)
}

// UpdateCategory godoc
//
//	@Summary		Update a category
//	@Description	Update a category by its ID
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Category ID"
//	@Param			request	body		dtos.UpdateCategoryRequest	true	"Category information"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/categories/{id} [patch]
func (cc *CategoryController) UpdateCategory(ctx *gin.Context) {
	categoryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	var categoryRequest dtos.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&categoryRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid request body").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	currentUserId, ok := ctx.Get("user_id")
	if !ok {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot get the current user's ID").
			Respond(ctx)
		return
	}

	reqCtx := ctx.Request.Context()
	if err := cc.categoryService.UpdateCategory(reqCtx, categoryId, categoryRequest, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to update the category").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	headers := headers.NewHeaders(nil, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully updated category").
		Respond(ctx)
}

// DeleteCategory godoc
//
//	@Summary		Delete a category
//	@Description	Delete a category by its ID
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Category ID"
//	@Success		204	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/categories/{id} [delete]
func (cc *CategoryController) DeleteCategory(ctx *gin.Context) {
	categoryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	currentUserId, ok := ctx.Get("user_id")
	if !ok {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot get the current user's ID").
			Respond(ctx)
		return
	}

	reqCtx := ctx.Request.Context()
	if err := cc.categoryService.DeleteCategory(reqCtx, categoryId, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to delete the category").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	headers := headers.NewHeaders(nil, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusNoContent).
		Respond(ctx)
}
