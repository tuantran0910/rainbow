package response

import "github.com/gin-gonic/gin"

type IAPIResponse interface {
	SetMessage(message string) *APIResponse
	SetHeaders(headers interface{}) *APIResponse
	SetData(data interface{}) *APIResponse
	SetPagination(pagination interface{}) *APIResponse
	SetStatusCode(status_code int) *APIResponse
	SetError(error string) *APIResponse
	Respond(ctx *gin.Context)
}

type APIResponse struct {
	Message    string      `json:"message"`
	Headers    interface{} `json:"headers"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
	StatusCode int         `json:"status_code"`
	Error      interface{} `json:"error"`
}

func NewAPIResponse() IAPIResponse {
	return &APIResponse{
		Message: "Success",
	}
}

func (r *APIResponse) SetMessage(message string) *APIResponse {
	r.Message = message
	return r
}

func (r *APIResponse) SetHeaders(headers interface{}) *APIResponse {
	r.Headers = headers
	return r
}

func (r *APIResponse) SetData(data interface{}) *APIResponse {
	r.Data = data
	return r
}

func (r *APIResponse) SetPagination(pagination interface{}) *APIResponse {
	r.Pagination = pagination
	return r
}

func (r *APIResponse) SetStatusCode(status_code int) *APIResponse {
	r.StatusCode = status_code
	return r
}

func (r *APIResponse) SetError(error string) *APIResponse {
	r.Error = error
	return r
}

func (r *APIResponse) Respond(ctx *gin.Context) {
	ctx.JSON(r.StatusCode, r)
}
