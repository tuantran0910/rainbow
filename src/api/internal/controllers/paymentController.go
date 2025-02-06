package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/services"
	"github.com/tuantran0910/rainbow/pkg/headers"
	"github.com/tuantran0910/rainbow/pkg/utils/response"
)

type PaymentController struct {
	paymentService services.IPaymentService
}

func NewPaymentController(paymentService services.IPaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// GetPayments godoc
//
//	@Summary		Get a list of payments
//	@Description	Get a list of payments
//	@Tags			Payment
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/payments [get]
func (pc *PaymentController) GetPayments(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	payments, err := pc.paymentService.GetPayments(reqCtx)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to get payments").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	paymentResponses := make([]*dtos.GetPaymentResponse, 0)
	for _, payment := range payments {
		paymentResponses = append(paymentResponses, &dtos.GetPaymentResponse{
			ID:        payment.ID,
			Method:    payment.Method,
			CreatedAt: payment.CreatedAt,
			UpdatedAt: payment.UpdatedAt,
			DeletedAt: payment.DeletedAt,
		})
	}
	data := &dtos.ListPaymentsResponse{
		Payments: paymentResponses,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved payments").
		SetData(data).
		Respond(ctx)
}
