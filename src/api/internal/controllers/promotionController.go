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
//	@Description	Get all promotions that are currently active (within their start and end date range)
//	@Tags			Promotion
//	@Produce		json
//	@Success		200	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/promotions [get]
func (pc *PromotionController) GetAllPromotions(ctx *gin.Context) {
	responseHeaders := headers.NewHeaders(nil, ctx)

	promotions, err := pc.promotionService.GetAllPromotions(ctx.Request.Context())
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

	response.
		NewAPIResponse().
		SetHeaders(responseHeaders).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully fetched promotions").
		SetData(promotions).
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
	if err := pc.promotionService.CreatePromotion(reqCtx, promotionRequest, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetHeaders(responseHeaders).
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to create promotion").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	response.NewAPIResponse().
		SetHeaders(responseHeaders).
		SetStatusCode(http.StatusCreated).
		SetMessage("Successfully created a promotion").
		Respond(ctx)
}
