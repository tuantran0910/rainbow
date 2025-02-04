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

type SellerController struct {
	sellerService services.ISellerService
}

func NewSellerController(sellerService services.ISellerService) *SellerController {
	return &SellerController{
		sellerService: sellerService,
	}
}

// GetSellers godoc
//
//	@Summary		Get a list of sellers
//	@Description	Get a list of sellers
//	@Tags			Seller
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"
//	@Param			limit	query		int	false	"Limit number"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/sellers [get]
func (sc *SellerController) GetSellers(ctx *gin.Context) {
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

	// Get a list of sellers
	sellers, pagination, err := sc.sellerService.GetSellers(reqCtx, page, limit)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch a list of sellers").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Parse output
	sellerResponses := make([]*dtos.GetSellerResponse, 0)
	for _, seller := range sellers {
		sellerResponses = append(sellerResponses, &dtos.GetSellerResponse{
			ID:        seller.ID,
			UserID:    seller.UserID,
			Name:      seller.Name,
			Link:      seller.Link,
			Logo:      seller.Logo,
			CreatedAt: seller.CreatedAt,
			UpdatedAt: seller.UpdatedAt,
			DeletedAt: seller.DeletedAt,
		})
	}
	data := &dtos.ListSellersResponse{
		Sellers: sellerResponses,
	}

	// Set headers
	headers := headers.NewHeaders(data, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetData(data).
		SetPagination(pagination).
		SetMessage("Successfully retrieved sellers").
		Respond(ctx)
}

// GetSellerById godoc
//
//	@Summary		Get a seller
//	@Description	Get a seller by its ID
//	@Tags			Seller
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Seller ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/sellers/{id} [get]
func (sc *SellerController) GetSellerById(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get ID from the URL
	sellerId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Get the seller
	seller, err := sc.sellerService.GetSellerById(reqCtx, sellerId)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch the seller").
			SetError(err.Error()).Respond(ctx)
		return
	}

	if seller == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("Seller not found").
			Respond(ctx)
		return
	}

	// Parse output
	data := &dtos.GetSellerResponse{
		ID:        seller.ID,
		UserID:    seller.UserID,
		Name:      seller.Name,
		Link:      seller.Link,
		Logo:      seller.Logo,
		CreatedAt: seller.CreatedAt,
		UpdatedAt: seller.UpdatedAt,
		DeletedAt: seller.DeletedAt,
	}

	// Set headers
	headers := headers.NewHeaders(data, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved seller").
		SetData(data).
		Respond(ctx)
}

// CreateSeller godoc
//
//	@Summary		Create Seller
//	@Description	Create a seller
//	@Tags			Seller
//	@Accept			json
//	@Produce		json
//	@Param			req	body		dtos.CreateSellerRequest	true	"Create Seller Request"
//	@Success		201	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/sellers [post]
func (sc *SellerController) CreateSeller(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Bind the request body to the struct
	var sellerRequest dtos.CreateSellerRequest
	if err := ctx.ShouldBindJSON(&sellerRequest); err != nil {
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

	// Create a seller
	if err := sc.sellerService.CreateSeller(reqCtx, sellerRequest, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to create a seller").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Set headers
	headers := headers.NewHeaders(nil, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusCreated).
		SetMessage("Successfully created a seller").
		Respond(ctx)
}

// UpdateSeller godoc
//
//	@Summary		Update Seller
//	@Description	Update a seller by its ID
//	@Tags			Seller
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Seller ID"
//	@Param			req	body		dtos.UpdateSellerRequest	true	"Update Seller Request"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/sellers/{id} [patch]
func (sc *SellerController) UpdateSeller(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get ID from the URL
	sellerId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Bind the request body to the struct
	var sellerRequest dtos.UpdateSellerRequest
	if err := ctx.ShouldBindJSON(&sellerRequest); err != nil {
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

	// Update the seller
	if err := sc.sellerService.UpdateSeller(reqCtx, sellerId, sellerRequest, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to update the seller").
			SetError(err.Error()).Respond(ctx)
		return
	}

	// Set headers
	headers := headers.NewHeaders(nil, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully updated the seller").
		Respond(ctx)
}

// DeleteSeller godoc
//
//	@Summary		Delete Seller
//	@Description	Delete a seller by its ID
//	@Tags			Seller
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Seller ID"
//	@Success		204	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/sellers/{id} [delete]
func (sc *SellerController) DeleteSeller(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get ID from the URL
	sellerId, err := uuid.Parse(ctx.Param("id"))
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

	// Delete the seller
	if err := sc.sellerService.DeleteSeller(reqCtx, sellerId, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to delete the seller").
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
