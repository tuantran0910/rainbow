package headers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Headers struct {
	ContentType   string `json:"content-type"   default:"application/json"`
	ContentLength int    `json:"content-length"`
	Date          string `json:"date"`
	Server        string `json:"server"`
	Endpoint      string `json:"endpoint"`
}

func NewHeaders(data interface{}, ctx *gin.Context) *Headers {
	// Calculate the content length
	jsonData, err := json.Marshal(data)
	if err != nil {
		jsonData = []byte{}
	}
	contentLength := len(jsonData)

	// Get other fields's value
	date := time.Now().Format(time.RFC1123)
	server := strings.Split(ctx.Request.Host, ":")[0]
	endpoint := ctx.Request.URL.String()

	return &Headers{
		ContentLength: contentLength,
		Date:          date,
		Server:        server,
		Endpoint:      endpoint,
	}
}
