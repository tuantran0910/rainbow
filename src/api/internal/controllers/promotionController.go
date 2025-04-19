package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/services"
	"github.com/tuantran0910/rainbow/pkg/headers"
	"github.com/tuantran0910/rainbow/pkg/utils/response"
)

type PromotionController struct {
	promotionService services.IPromotionService
}

func NewPromotionController(promotionService services.IPromotionService) *PromotionController {
	return &PromotionController{
		promotionService: promotionService,
	}
}

// GetAllPromotions godoc
//
//	@Summary		Get all active promotions
//	@Description	Get all promotions that are currently active (within their start and end date range). If user_id is provided, returns only promotions available for that user.
//	@Tags			Promotion
//	@Produce		json
//	@Param			user_id	query		string	false	"User ID to filter available promotions"
//	@Success		200		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/promotions [get]
func (pc *PromotionController) GetAllPromotions(ctx *gin.Context) {
	responseHeaders := headers.NewHeaders(nil, ctx)

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
	promotions, err := pc.promotionService.GetAllPromotions(reqCtx, currentUserId.(uuid.UUID))
	if err != nil {
		response.
			NewAPIResponse().
			SetHeaders(responseHeaders).
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch promotions").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	promotionResponses := make([]*dtos.GetPromotionResponse, 0)
	for _, promotion := range promotions {
		promotionResponses = append(promotionResponses, &dtos.GetPromotionResponse{
			ID:            promotion.ID,
			Name:          promotion.Name,
			DiscountType:  promotion.DiscountType,
			DiscountValue: promotion.DiscountValue,
			StartDate:     promotion.StartDate,
			EndDate:       promotion.EndDate,
			MaxUses:       promotion.MaxUses,
			UsedCount:     promotion.UsedCount,
			CreatedAt:     promotion.CreatedAt,
			UpdatedAt:     promotion.UpdatedAt,
			DeletedAt:     promotion.DeletedAt,
		})
	}
	data := &dtos.ListPromotionsResponse{
		Promotions: promotionResponses,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully fetched promotions").
		SetData(data).
		Respond(ctx)
}

// CreatePromotion godoc
//
//	@Summary		Create a promotion
//	@Description	Create a promotion
//	@Tags			Promotion
//	@Accept			json
//	@Produce		json
//	@Param			req	body		dtos.CreatePromotionRequest	true	"Create Promotion Request"
//	@Success		201	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/promotions [post]
func (pc *PromotionController) CreatePromotion(ctx *gin.Context) {
	responseHeaders := headers.NewHeaders(nil, ctx)

	var promotionRequest dtos.CreatePromotionRequest
	if err := ctx.ShouldBindJSON(&promotionRequest); err != nil {
		response.
			NewAPIResponse().
			SetHeaders(responseHeaders).
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
	promotion, err := pc.promotionService.CreatePromotion(
		reqCtx,
		promotionRequest,
		currentUserId.(uuid.UUID),
	)
	if err != nil {
		response.
			NewAPIResponse().
			SetHeaders(responseHeaders).
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to create promotion").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	data := &dtos.GetPromotionResponse{
		ID:            promotion.ID,
		Name:          promotion.Name,
		DiscountType:  promotion.DiscountType,
		DiscountValue: promotion.DiscountValue,
		StartDate:     promotion.StartDate,
		EndDate:       promotion.EndDate,
		MaxUses:       promotion.MaxUses,
		UsedCount:     promotion.UsedCount,
		CreatedAt:     promotion.CreatedAt,
		UpdatedAt:     promotion.UpdatedAt,
	}

	responseHeaders = headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(responseHeaders).
		SetStatusCode(http.StatusCreated).
		SetMessage("Successfully created a promotion").
		SetData(data).
		Respond(ctx)
}
