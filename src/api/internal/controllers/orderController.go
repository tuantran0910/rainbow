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

type OrderController struct {
	orderService services.IOrderService
}

func NewOrderController(orderService services.IOrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (oc *OrderController) GetOrdersByUserId(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid pagination's parameter page, page must be a positive integer").
			SetError(err.Error()).
			Respond(ctx)
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
	orders, pagination, err := oc.orderService.GetOrdersByUserId(
		reqCtx,
		page,
		limit,
		currentUserId.(uuid.UUID),
	)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch user's orders").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	orderResponses := make([]*dtos.GetOrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, &dtos.GetOrderResponse{
			ID:              order.ID,
			UserID:          order.UserID,
			PaymentID:       order.PaymentID,
			ShippingAddress: order.ShippingAddress,
			TotalAmount:     order.TotalAmount,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
			DeletedAt:       order.DeletedAt,
		})
	}
	data := dtos.ListOrdersResponse{
		Orders: orderResponses,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetData(data).
		SetPagination(pagination).
		SetMessage("Successfully retrieved user's orders").
		Respond(ctx)
}

func (oc *OrderController) GetOrderById(ctx *gin.Context) {
	orderId, err := uuid.Parse(ctx.Param("id"))
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
	order, err := oc.orderService.GetOrderById(reqCtx, orderId, currentUserId.(uuid.UUID))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch the order").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	if order == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("Order not found").
			Respond(ctx)
		return
	}

	orderItems := make([]*dtos.GetOrderItemResponse, 0, len(order.OrderItems))
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, &dtos.GetOrderItemResponse{
			ID:        orderItem.ID,
			OrderID:   orderItem.OrderID,
			BookID:    orderItem.BookID,
			Quantity:  orderItem.Quantity,
			UnitPrice: orderItem.UnitPrice,
			Discount:  orderItem.Discount,
			CreatedAt: orderItem.CreatedAt,
			UpdatedAt: orderItem.UpdatedAt,
			DeletedAt: orderItem.DeletedAt,
		})
	}
	data := &dtos.GetOrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		PaymentID:       order.PaymentID,
		ShippingAddress: order.ShippingAddress,
		TotalAmount:     order.TotalAmount,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		DeletedAt:       order.DeletedAt,
		OrderItems:      &orderItems,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetData(data).
		SetMessage("Successfully retrieved order").
		Respond(ctx)
}

func (oc *OrderController) CreateOrder(ctx *gin.Context) {
	var orderRequest dtos.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid request's body").
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
	if err := oc.orderService.CreateOrder(reqCtx, orderRequest, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to create order").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	headers := headers.NewHeaders(nil, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusCreated).
		SetMessage("Successfully created order").
		Respond(ctx)
}

func (oc *OrderController) DeleteOrder(ctx *gin.Context) {
	orderId, err := uuid.Parse(ctx.Param("id"))
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
	if err := oc.orderService.DeleteOrder(reqCtx, orderId, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to delete order").
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
