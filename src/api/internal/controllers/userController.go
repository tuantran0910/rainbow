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

type UserController struct {
	userService services.IUserService
}

func NewUserController(userService services.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetUsers godoc
//
//	@Summary		Get a list of users
//	@Description	Get a list of users
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page"
//	@Param			limit	query		int	false	"Limit"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/users [get]
func (uc *UserController) GetUsers(ctx *gin.Context) {
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

	// Get a list of users
	users, pagination, err := uc.userService.GetUsers(reqCtx, page, limit)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch a list of users").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Parse output
	userResponses := make([]*dtos.GetUserResponse, 0)
	for _, user := range users {
		userResponses = append(userResponses, &dtos.GetUserResponse{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			LastLogin:   user.LastLogin,
			IsActive:    user.IsActive,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			DeletedAt:   user.DeletedAt,
		})
	}
	data := &dtos.ListUsersResponse{
		Users: userResponses,
	}

	// Set headers
	headers := headers.NewHeaders(data, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetPagination(pagination).
		SetMessage("Successfully retrieved users").
		SetData(data).
		Respond(ctx)
}

// GetCurrentUser godoc
//
//	@Summary		Get the current user
//	@Description	Get the current user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/users/me [get]
func (uc *UserController) GetCurrentUser(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get the current user's ID
	currentUserId, ok := ctx.Get("user_id")
	if !ok {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot get the current user's ID").
			Respond(ctx)
		return
	}

	// Get a user by ID
	user, err := uc.userService.GetUserById(reqCtx, currentUserId.(uuid.UUID))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch a user").
			SetError(err.Error()).Respond(ctx)
		return
	}

	if user == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("User not found").
			Respond(ctx)
		return
	}

	// Parse output
	data := &dtos.GetUserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		LastLogin:   user.LastLogin,
		IsActive:    user.IsActive,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}

	// Set headers
	headers := headers.NewHeaders(data, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved user").
		SetData(data).
		Respond(ctx)
}

// GetUserById godoc
//
//	@Summary		Get a user by ID
//	@Description	Get a user by ID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/users/{id} [get]
func (uc *UserController) GetUserById(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get ID from the URL
	userId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Get a user by ID
	user, err := uc.userService.GetUserById(reqCtx, userId)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch a user").
			SetError(err.Error()).Respond(ctx)
		return
	}

	if user == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("User not found").
			Respond(ctx)
		return
	}

	// Parse output
	data := &dtos.GetUserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		LastLogin:   user.LastLogin,
		IsActive:    user.IsActive,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}

	// Set headers
	headers := headers.NewHeaders(data, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved user").
		SetData(data).
		Respond(ctx)
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update a user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string					true	"User ID"
//	@Param			req	body		dtos.UpdateUserRequest	true	"Update User Request"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/users/{id} [patch]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get ID from the URL
	userId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Get the request body
	var userRequest dtos.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid request body").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Get the current user's ID
	currentUserId, ok := ctx.Get("user_id")
	if !ok {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot get the current user's ID").
			Respond(ctx)
		return
	}

	// Update a user
	if err := uc.userService.UpdateUser(reqCtx, userId, userRequest, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to update a user").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Set headers
	headers := headers.NewHeaders(nil, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully updated user").
		Respond(ctx)
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		204	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/users/{id} [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get ID from the URL
	userId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Get the current user's ID
	currentUserId, ok := ctx.Get("user_id")
	if !ok {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot get the current user's ID").
			Respond(ctx)
		return
	}

	// Delete a user
	if err := uc.userService.DeleteUser(reqCtx, userId, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to delete a user").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Set headers
	headers := headers.NewHeaders(nil, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusNoContent).
		Respond(ctx)
}
