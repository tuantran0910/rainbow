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

type AuthorController struct {
	authorService services.IAuthorService
}

func NewAuthorController(authorService services.IAuthorService) *AuthorController {
	return &AuthorController{
		authorService: authorService,
	}
}

// GetAuthors godoc
//
//	@Summary		Get Authors
//	@Description	Fetch a list of authors
//	@Tags			Author
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"
//	@Param			limit	query		int	false	"Number of items per page"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/authors [get]
func (ac *AuthorController) GetAuthors(ctx *gin.Context) {
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
	authors, pagination, err := ac.authorService.GetAuthors(reqCtx, page, limit)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to get authors").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	authorResponses := make([]*dtos.GetAuthorResponse, 0)
	for _, author := range authors {
		authorResponses = append(authorResponses, &dtos.GetAuthorResponse{
			ID:        author.ID,
			Name:      author.Name,
			Slug:      author.Slug,
			CreatedAt: author.CreatedAt,
			UpdatedAt: author.UpdatedAt,
			DeletedAt: author.DeletedAt,
		})
	}
	data := &dtos.ListAuthorsResponse{
		Authors: authorResponses,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetPagination(pagination).
		SetMessage("Successfully retrieved authors").
		SetData(data).
		Respond(ctx)
}

// GetAuthorById godoc
//
//	@Summary		Get Author by ID
//	@Description	Fetch an author by ID
//	@Tags			Author
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Author ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/authors/{id} [get]
func (ac *AuthorController) GetAuthorById(ctx *gin.Context) {
	authorId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
		return
	}

	reqCtx := ctx.Request.Context()
	author, err := ac.authorService.GetAuthorById(reqCtx, authorId)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to get the author").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	if author == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("Author not found").
			Respond(ctx)
		return
	}

	data := &dtos.GetAuthorResponse{
		ID:        author.ID,
		Name:      author.Name,
		Slug:      author.Slug,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
		DeletedAt: author.DeletedAt,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved author").
		SetData(data).
		Respond(ctx)
}

// GetAuthorBySlug godoc
//
//	@Summary		Get Author by Slug
//	@Description	Fetch an author by its slug
//	@Tags			Author
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"Author Slug"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/authors/slug/{slug} [get]
func (ac *AuthorController) GetAuthorBySlug(ctx *gin.Context) {
	authorSlug := ctx.Param("slug")

	reqCtx := ctx.Request.Context()
	author, err := ac.authorService.GetAuthorBySlug(reqCtx, authorSlug)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to get the author").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	if author == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("Author not found").
			Respond(ctx)
		return
	}

	data := &dtos.GetAuthorResponse{
		ID:        author.ID,
		Name:      author.Name,
		Slug:      author.Slug,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
		DeletedAt: author.DeletedAt,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved author").
		SetData(data).
		Respond(ctx)
}

// CreateAuthor godoc
//
//	@Summary		Create Author
//	@Description	Create a new author
//	@Tags			Author
//	@Accept			json
//	@Produce		json
//	@Param			req	body		dtos.CreateAuthorRequest	true	"Create Author Request"
//	@Success		201	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/authors [post]
func (ac *AuthorController) CreateAuthor(ctx *gin.Context) {
	var authorRequest dtos.CreateAuthorRequest
	if err := ctx.ShouldBindJSON(&authorRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
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
	err := ac.authorService.CreateAuthor(reqCtx, authorRequest, currentUserId.(uuid.UUID))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to create the author").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	headers := headers.NewHeaders(nil, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusCreated).
		SetMessage("Successfully created author").
		Respond(ctx)
}

// UpdateAuthor godoc
//
//	@Summary		Update Author
//	@Description	Update an author by its ID
//	@Tags			Author
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Author ID"
//	@Param			req	body		dtos.UpdateAuthorRequest	true	"Update Author Request"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/authors/{id} [patch]
func (ac *AuthorController) UpdateAuthor(ctx *gin.Context) {
	authorId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
		return
	}

	var authorRequest dtos.UpdateAuthorRequest
	if err := ctx.ShouldBindJSON(&authorRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
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
	err = ac.authorService.UpdateAuthor(reqCtx, authorId, authorRequest, currentUserId.(uuid.UUID))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to update the author").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	headers := headers.NewHeaders(nil, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully updated author").
		Respond(ctx)
}

// DeleteAuthor godoc
//
//	@Summary		Delete Author
//	@Description	Delete an author by its ID
//	@Tags			Author
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Author ID"
//	@Success		204	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/authors/{id} [delete]
func (ac *AuthorController) DeleteAuthor(ctx *gin.Context) {
	authorId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
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
	err = ac.authorService.DeleteAuthor(reqCtx, authorId, currentUserId.(uuid.UUID))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to delete the author").
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
