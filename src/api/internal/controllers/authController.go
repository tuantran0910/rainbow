package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/services"
	"github.com/tuantran0910/rainbow/pkg/headers"
	"github.com/tuantran0910/rainbow/pkg/utils/response"
)

type AuthController struct {
	authService services.IAuthService
}

func NewAuthController(authService services.IAuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Login(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get the body request
	var loginUserRequest dtos.LoginUserRequest
	if err := ctx.ShouldBindJSON(&loginUserRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Get the token
	token, err := ac.authService.Login(reqCtx, loginUserRequest)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusUnauthorized).
			SetMessage("Unauthorized").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Set headers
	headers := headers.NewHeaders(nil, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetData(token).
		SetMessage("Successfully logged in").
		Respond(ctx)
}

func (ac *AuthController) Register(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get the body request
	var registerUserRequest dtos.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&registerUserRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Create a user
	if err := ac.authService.Register(reqCtx, registerUserRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Cannot create user").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Set headers
	headers := headers.NewHeaders(nil, ctx)

	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusCreated).
		SetMessage("Successfully created user").
		Respond(ctx)
}
