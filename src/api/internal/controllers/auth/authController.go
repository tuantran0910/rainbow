package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/tuantran0910/rainbow/internal/models/user"
	service "github.com/tuantran0910/rainbow/internal/services/auth"
	"github.com/tuantran0910/rainbow/pkg/headers"
	"github.com/tuantran0910/rainbow/pkg/utils/response"
)

type AuthController struct {
	authService service.IAuthService
}

func NewAuthController(authService service.IAuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Login(ctx *gin.Context) {
	// Get the context
	reqCtx := ctx.Request.Context()

	// Get the body request
	var userRequest model.UserRequest
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Get the token
	token, err := ac.authService.Login(reqCtx, *userRequest.Email, *userRequest.Password)
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
	var userRequest model.UserRequest
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	// Create a user
	if err := ac.authService.Register(reqCtx, userRequest); err != nil {
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
