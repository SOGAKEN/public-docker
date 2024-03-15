package handlers

import (
	"github.com/gin-gonic/gin"
)

// APIHandler defines the interface for handling requests to different APIs.
type APIHandler interface {
	HandleRequest(c *gin.Context, requestData interface{}) (interface{}, error)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
