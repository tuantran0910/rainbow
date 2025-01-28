package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuantran0910/rainbow/pkg/utils/response"
)

type BaseController struct{}

func NewBaseController() (*BaseController, error) {
	return &BaseController{}, nil
}

func (bc *BaseController) HealthCheck(ctx *gin.Context) {
	response.NewAPIResponse().SetStatusCode(http.StatusOK).SetMessage("ok!").Respond(ctx)
}

func (bc *BaseController) HomePage(ctx *gin.Context) {
	response.NewAPIResponse().SetStatusCode(http.StatusOK).SetData(
		map[string]string{
			"author":      "tntuan0910@gmail.com",
			"github":      "https://github.com/tuantran0910/rainbow",
			"description": "Simple Backend API for Rainbow Data Platform",
		},
	).Respond(ctx)
}
